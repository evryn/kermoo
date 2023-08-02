package user_config

import (
	"buggybox/modules/planner"
	"buggybox/modules/process"
	"buggybox/modules/web_server"
)

var Prepared PreparedConfigType

type PreparedConfigType struct {
	SchemaVersion string
	Process       process.Process
	Plans         []*planner.Plan
	WebServers    []web_server.WebServer
}

func GetPreparedConfig(u UserConfigType) PreparedConfigType {
	// Preserve process killer
	prepared := PreparedConfigType{
		Process:    u.Process,
		WebServers: u.WebServers,
	}

	// Preserve global plans
	for i := range u.Plans {
		prepared.Plans = append(prepared.Plans, &u.Plans[i])
	}

	u.Process.GetDesiredPlanNames()

	plannable := &u.Process
	prepared.preparePlannable(plannable)

	for _, ws := range prepared.WebServers {
		plannable := &ws
		prepared.preparePlannable(plannable)
	}

	return prepared
}

func (u *PreparedConfigType) preparePlannable(plannable planner.Plannable) {
	desiredPlans := []string{}

	if plannable.HasCustomPlan() {
		planName := plannable.GetUid()
		customPlan := plannable.MakeCustomPlan()
		customPlan.Name = &planName
		u.Plans = append(u.Plans, customPlan)
		desiredPlans = append(desiredPlans, planName)
	} else {
		desiredPlans = plannable.GetDesiredPlanNames()
	}

	plansToAssign := []*planner.Plan{}
	for _, planName := range desiredPlans {
		plansToAssign = append(plansToAssign, u.findPlan(planName))
	}

	plannable.AssignPlans(plansToAssign)
}

func (u *PreparedConfigType) findPlan(name string) *planner.Plan {
	for _, plan := range u.Plans {
		if plan.Name != nil && *plan.Name == name {
			return plan
		}
	}

	return nil
}

func (u *PreparedConfigType) findDuplicatePlans() string {
	for _, v := range u.Plans {

	}
}
