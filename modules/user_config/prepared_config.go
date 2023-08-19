package user_config

import (
	"fmt"
	"kermoo/modules/cpu"
	"kermoo/modules/logger"
	"kermoo/modules/memory"
	"kermoo/modules/planner"
	"kermoo/modules/process"
	"kermoo/modules/utils"
	"kermoo/modules/web_server"
	"strings"
	"time"

	"go.uber.org/zap"
)

var Prepared PreparedConfigType

type PreparedConfigType struct {
	SchemaVersion string
	Process       *process.Process
	Cpu           *cpu.Cpu
	Memory        *memory.Memory
	Plans         []*planner.Plan
	WebServers    []*web_server.WebServer
}

func (pc *PreparedConfigType) Start() {
	if pc.Process != nil && pc.Process.Delay != nil {
		dur, _ := pc.Process.Delay.ToStandardDuration()
		logger.Log.Info("sleeping because of process manager configuration...", zap.Duration("sleep", dur))
		time.Sleep(dur)
		logger.Log.Info("woke up.")
	}

	for _, plan := range pc.Plans {
		go plan.Start()
	}
}

func (u *PreparedConfigType) preparePlannable(plannable planner.Plannable) error {
	desiredPlans := plannable.GetDesiredPlanNames()

	if len(desiredPlans) == 0 {
		planName := plannable.GetName()
		var dedicatedPlan *planner.Plan

		if plannable.HasInlinePlan() {
			dedicatedPlan = plannable.MakeInlinePlan()
			planName += "-custom-plan"
		} else {
			dedicatedPlan = plannable.MakeDefaultPlan()
			planName += "-default-plan"
		}

		dedicatedPlan.Name = &planName
		dedicatedPlan.MakePrivate()

		u.Plans = append(u.Plans, dedicatedPlan)
		desiredPlans = append(desiredPlans, planName)
	}

	logger.Log.Debug("preparing plannable", zap.String("plannable", plannable.GetName()), zap.Any("desired_plans", desiredPlans))

	for _, planName := range desiredPlans {
		desiredPlan := u.findPlan(planName)

		if desiredPlan == nil {
			return fmt.Errorf("plan %s not found for sub-app %s", planName, plannable.GetName())
		}

		desiredPlan.Assign(plannable)
	}

	return nil
}

func (u *PreparedConfigType) findPlan(name string) *planner.Plan {
	for _, plan := range u.Plans {
		if plan.Name != nil && *plan.Name == name {
			return plan
		}
	}

	return nil
}

func (pc *PreparedConfigType) validateDuplicateApps() error {
	dups := pc.findDuplicateApps()
	if len(dups) > 0 {
		return fmt.Errorf("there are duplicate sub-apps: %s", strings.Join(dups, ", "))
	}

	return nil
}

func (pc *PreparedConfigType) validateDuplicatePlans() error {
	dups := pc.findDuplicatePlans()
	if len(dups) > 0 {
		return fmt.Errorf("there are duplicate plans: %s", strings.Join(dups, ", "))
	}

	return nil
}

func (pc *PreparedConfigType) validatePlans() error {
	for _, plan := range pc.Plans {
		err := plan.Validate()
		if err != nil {
			return fmt.Errorf("plan %s is invalid: %v", *plan.Name, err)
		}
	}

	return nil
}

func (pc *PreparedConfigType) validateProcess() error {
	if pc.Process == nil {
		return nil
	}

	if err := pc.Process.Validate(); err != nil {
		return fmt.Errorf("process manager is invalid: %v", err)
	}

	return nil
}

func (pc *PreparedConfigType) validateCpu() error {
	if pc.Cpu == nil {
		return nil
	}

	if err := pc.Cpu.Validate(); err != nil {
		return fmt.Errorf("cpu manager is invalid: %v", err)
	}

	return nil
}

func (pc *PreparedConfigType) validateMemory() error {
	if pc.Memory == nil {
		return nil
	}

	if err := pc.Memory.Leak.Validate(); err != nil {
		return fmt.Errorf("memory leaker is invalid: %v", err)
	}

	return nil
}

func (pc *PreparedConfigType) validateWebservers() error {
	for _, webServer := range pc.WebServers {
		err := webServer.Validate()
		if err != nil {
			return fmt.Errorf("webserver %s is invalid: %v", webServer.GetName(), err)
		}

		for _, route := range webServer.Routes {
			err := route.Validate()

			if err != nil {
				return fmt.Errorf("route %s is invalid for webserver %s: %v", route.Path, webServer.GetName(), err)
			}
		}
	}

	return nil
}

func (pc *PreparedConfigType) Validate() error {
	if err := pc.validateDuplicateApps(); err != nil {
		return err
	}

	if err := pc.validateDuplicatePlans(); err != nil {
		return err
	}

	if err := pc.validatePlans(); err != nil {
		return err
	}

	if err := pc.validateProcess(); err != nil {
		return err
	}

	if err := pc.validateCpu(); err != nil {
		return err
	}

	if err := pc.validateMemory(); err != nil {
		return err
	}

	if err := pc.validateWebservers(); err != nil {
		return err
	}

	return nil
}

func (u *PreparedConfigType) findDuplicateApps() []string {
	apps := []string{
		u.Process.GetName(),
	}

	for _, v := range u.WebServers {
		apps = append(apps, v.GetName())
	}

	return utils.GetDuplicates(apps)
}

func (u *PreparedConfigType) findDuplicatePlans() []string {
	plans := []string{}

	for _, v := range u.Plans {
		plans = append(plans, *v.Name)
	}

	return utils.GetDuplicates(plans)
}
