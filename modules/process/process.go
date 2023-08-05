package process

import (
	"buggybox/modules/common"
	"buggybox/modules/logger"
	"buggybox/modules/planner"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

// Ensure that Process is implementing Plannable
var _ planner.Plannable = &Process{}

type Process struct {
	planner.PlannableTrait
	Exit struct {
		After common.SingleValueDur `json:"after"`
		Code  uint                  `json:"code"`
	} `json:"exit"`
}

func (p *Process) MustRun() {
	plan := p.MakeCustomPlan()
	err := plan.ExecuteAll()
	if err != nil {
		logger.Log.Fatal("process manager plan execution failed", zap.Error(err))
	}
}

func (p *Process) GetUid() string {
	return "process-manager"
}

func (p *Process) HasCustomPlan() bool {
	return true
}

func (p Process) GetDesiredPlanNames() []string {
	return nil
}

func (p Process) ShouldStart() bool {
	value, _ := p.Exit.After.GetValue()

	return value != 0
}

func (p Process) Validate() error {
	_, err := p.Exit.After.GetValue()

	if err != nil {
		return fmt.Errorf("unable to get exit duration: %v", err)
	}

	return nil
}

func (p *Process) GetPlanCallbacks() planner.Callbacks {
	return planner.Callbacks{
		PreSleep: func(ep *planner.ExecutablePlan, ev *planner.ExecutableValue) planner.PlanSignal {
			return planner.PLAN_SIGNAL_CONTINUE
		},
		PostSleep: func(startedAt time.Time, timeSpent time.Duration) planner.PlanSignal {
			logger.Log.Info("process is exiting due to the specified alive time in configuration",
				zap.Duration("seconds_alive", timeSpent),
				zap.Int("exit_code", int(p.Exit.Code)),
			)

			os.Exit(int(p.Exit.Code))

			return planner.PLAN_SIGNAL_TERMINATE
		},
	}
}

func (p *Process) MakeCustomPlan() *planner.Plan {
	value, _ := p.Exit.After.GetValue()
	name := p.GetUid()

	plan := planner.InitPlan(planner.Plan{
		Name:     &name,
		Interval: &value,
		Duration: &value,
	})

	// Set a dummy value since plan validation requires it
	dummyValue := float32(1.0)
	plan.Value.Exactly = &dummyValue

	return &plan
}
