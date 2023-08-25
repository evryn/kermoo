package process

import (
	"kermoo/modules/fluent"
	"kermoo/modules/logger"
	"kermoo/modules/planner"
	"os"

	"go.uber.org/zap"
)

// Ensure that Process is implementing Plannable
var _ planner.Plannable = &Process{}

type Process struct {
	planner.CanAssignPlan

	// Delay optionally defines the initial startup delay. The process will sleep until
	// that delay is reached.
	Delay *fluent.FluentDuration `json:"delay"`

	// Exit optionally simulates sudden termination of the process in the given time with
	// the given exit code.
	Exit *ProcessExit `json:"exit"`
}

type ProcessExit struct {
	// After determines the duration in which the process will be terminated
	After fluent.FluentDuration `json:"after"`

	// Code indicates the exit code of the process when the time is reached.
	Code uint `json:"code"`
}

func (p *Process) GetName() string {
	return "process-manager"
}

func (p *Process) HasInlinePlan() bool {
	return p.MakeInlinePlan() != nil
}

func (p Process) GetDesiredPlanNames() []string {
	return nil
}

func (p Process) Validate() error {
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
	if p.Exit == nil {
		return nil
	}

	name := p.GetName()

	valueDur := p.Exit.After

	plan := planner.NewPlan(planner.Plan{
		Name:     &name,
		Interval: &valueDur,
		// Duration: &valueDur,
	})

	// Set a dummy value since plan validation requires it
	plan.Percentage = fluent.NewMustFluentFloat("100")

	return &plan
}

func (p *Process) MakeDefaultPlan() *planner.Plan {
	return nil
}
