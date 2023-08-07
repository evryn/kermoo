package planner_test

import (
	"kermoo/modules/planner"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type PlanRecorder struct {
	planner.PlannableTrait
	TotalTimeSpent time.Duration
	Executions     []Execution
	ExecutaionCap  int
}

type Execution struct {
	TimeSpent time.Duration
	Cycles    planner.Cycle
}

type ExpectedExecutionValue struct {
	Static    float32
	Minimum   float32
	Maximum   float32
	IsBetween bool
}

func (r *PlanRecorder) Reset() {
	r.Executions = []Execution{}
	r.TotalTimeSpent = 0
}

func (r *PlanRecorder) AssertExpectedValues(t *testing.T, expectedValues []ExpectedExecutionValue) {
	assert.Len(t, r.Executions, len(expectedValues))

	for i, ev := range expectedValues {
		v := r.Executions[i].Cycles.ExecutableValue

		if ev.IsBetween {
			assert.GreaterOrEqual(t, v.GetValue(), ev.Minimum)
			assert.LessOrEqual(t, v.GetValue(), ev.Maximum)
			assert.Equal(t, v.GetValue(), v.GetValue(), "Ranged values should always get similar result which is determined by the first attempt to retrieve the value.")
		} else {
			assert.Equal(t, ev.Static, v.GetValue())
		}
	}
}

func (r *PlanRecorder) AssertTotalTimeSpent(t *testing.T, expectedDuration time.Duration, expectedError float64) {
	assert.LessOrEqual(
		t,
		r.TotalTimeSpent-expectedDuration,
		time.Duration(float64(expectedDuration)*expectedError),
	)
}

func (r *PlanRecorder) GetUid() string {
	return "recorder"
}

func (r *PlanRecorder) GetDesiredPlanNames() []string {
	return nil
}

func (r *PlanRecorder) HasCustomPlan() bool {
	return false
}

func (r *PlanRecorder) MakeCustomPlan() *planner.Plan {
	return nil
}

func (r *PlanRecorder) MakeDefaultPlan() *planner.Plan {
	return nil
}

func (r *PlanRecorder) GetPlanCycleHooks() planner.CycleHooks {
	preSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		r.Executions = append(r.Executions, Execution{
			Cycles: cycle,
		})

		return planner.PLAN_SIGNAL_CONTINUE
	})

	postSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		r.Executions[len(r.Executions)-1].TimeSpent = cycle.TimeSpent
		r.TotalTimeSpent += cycle.TimeSpent

		if r.ExecutaionCap != 0 && len(r.Executions) == r.ExecutaionCap {
			return planner.PLAN_SIGNAL_TERMINATE
		}

		return planner.PLAN_SIGNAL_CONTINUE
	})

	return planner.CycleHooks{
		PreSleep:  &preSleep,
		PostSleep: &postSleep,
	}
}

var Recorder PlanRecorder
