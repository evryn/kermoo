package process

import (
	"buggybox/modules/common"
	"buggybox/modules/logger"
	"buggybox/modules/planner"
	"os"
	"time"

	"go.uber.org/zap"
)

type Process struct {
	Exit struct {
		After common.SingleValueDur `json:"after"`
		Code  uint                  `json:"code"`
	} `json:"exit"`
}

func (p *Process) MustRun() {
	plan := p.makePlan()
	err := plan.ExecuteAll()
	if err != nil {
		logger.Log.Fatal("process killer plan execution failed", zap.Error(err))
	}
}

func (p *Process) makePlan() planner.Plan {
	value, _ := p.Exit.After.GetValue()
	name := "process-killer-plan"

	plan := planner.InitPlan(planner.Plan{
		Name:     &name,
		Interval: &value,
		Duration: &value,
	})

	// Set a dummy value since plan validation requires it
	dummyValue := float32(1.0)
	plan.Value.Exactly = &dummyValue

	plan.MakePrivate()

	callback := planner.Callbacks{
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

	plan.AddCallback(callback)

	return plan
}
