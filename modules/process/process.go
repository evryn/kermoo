package process

import (
	"fmt"
	"kermoo/modules/common"
	"kermoo/modules/logger"
	"kermoo/modules/planner"
	"kermoo/modules/utils"
	"os"

	"go.uber.org/zap"
)

// Ensure that Process is implementing Plannable
var _ planner.Plannable = &Process{}

type Process struct {
	planner.PlannableTrait
	Delay *common.SingleDuration `json:"delay"`
	Exit  *ProcessExit           `json:"exit"`
}

type ProcessExit struct {
	After common.SingleDuration `json:"after"`
	Code  uint                  `json:"code"`
}

func (p *Process) GetName() string {
	return "process-manager"
}

func (p *Process) HasInlinePlan() bool {
	return p.Exit != nil
}

func (p Process) GetDesiredPlanNames() []string {
	return nil
}

func (p Process) Validate() error {
	if p.Delay != nil {
		if _, err := p.Delay.ToStandardDuration(); err != nil {
			return fmt.Errorf("unable to get delay duration: %v", err)
		}
	}

	if p.Exit != nil {
		if _, err := p.Exit.After.ToStandardDuration(); err != nil {
			return fmt.Errorf("unable to get exit duration: %v", err)
		}
	}

	return nil
}

func (p *Process) GetPlanCycleHooks() planner.CycleHooks {
	preSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		return planner.PLAN_SIGNAL_CONTINUE
	})

	postSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		logger.Log.Info("process is exiting due to the specified alive time in configuration",
			zap.Duration("seconds_alive", cycle.TimeSpent),
			zap.Int("exit_code", int(p.Exit.Code)),
		)

		os.Exit(int(p.Exit.Code))

		return planner.PLAN_SIGNAL_TERMINATE
	})

	return planner.CycleHooks{
		PreSleep:  &preSleep,
		PostSleep: &postSleep,
	}
}

func (p *Process) MakeInlinePlan() *planner.Plan {
	value, _ := p.Exit.After.ToStandardDuration()
	name := p.GetName()

	valueDur := common.Duration(value)

	plan := planner.InitPlan(planner.Plan{
		Name:     &name,
		Interval: &valueDur,
		Duration: &valueDur,
	})

	// Set a dummy value since plan validation requires it
	plan.Value = &common.MultiFloat{
		SingleFloat: common.SingleFloat{
			Exactly: utils.NewP[float32](1.0),
		},
	}

	return &plan
}

func (p *Process) MakeDefaultPlan() *planner.Plan {
	return nil
}
