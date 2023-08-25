package planner_test

import (
	"kermoo/modules/fluent"
	"kermoo/modules/planner"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PlanRecorder struct {
	planner.CanAssignPlan
	TotalTimeSpent time.Duration
	Cycles         []planner.Cycle
	ExecutaionCap  int
}

type ExpectedCycleValue struct {
	Percentage *fluent.FluentFloat
	Size       *fluent.FluentSize
}

func (r *PlanRecorder) Reset() {
	r.Cycles = []planner.Cycle{}
	r.TotalTimeSpent = 0
}

func (r *PlanRecorder) AssertCycleValues(t *testing.T, expectedCycleValues []ExpectedCycleValue) {
	require.Len(t, r.Cycles, len(expectedCycleValues))

	for i, ev := range expectedCycleValues {
		if ev.Percentage != nil {
			actual := r.Cycles[i].Value.Percentage

			evpv := ev.Percentage.GetParsedValue()

			if evpv.IsRanged() {
				min, max, _ := ev.Percentage.GetParsedValue().GetRange()
				assert.Less(t, min, actual)
				assert.Greater(t, max, actual)

			} else {
				assert.Equal(t, ev.Percentage.Get(), actual)
			}

		}

		if ev.Size != nil {
			actual := r.Cycles[i].Value.Size

			evpv := ev.Size.GetParsedValue()

			if evpv.IsRanged() {
				max, min, _ := ev.Size.GetParsedValue().GetRange()
				assert.Less(t, min, actual)
				assert.Greater(t, max, actual)

			} else {
				assert.Equal(t, ev.Size.Get(), actual)
			}

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

func (r *PlanRecorder) GetName() string {
	return "recorder"
}

func (r *PlanRecorder) GetDesiredPlanNames() []string {
	return nil
}

func (r *PlanRecorder) HasInlinePlan() bool {
	return false
}

func (r *PlanRecorder) MakeInlinePlan() *planner.Plan {
	return nil
}

func (r *PlanRecorder) MakeDefaultPlan() *planner.Plan {
	return nil
}

func (r *PlanRecorder) GetPlanCycleHooks() planner.CycleHooks {
	preSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		r.Cycles = append(r.Cycles, cycle)

		return planner.PLAN_SIGNAL_CONTINUE
	})

	postSleep := planner.HookFunc(func(cycle planner.Cycle) planner.PlanSignal {
		r.Cycles[len(r.Cycles)-1].TimeSpent = cycle.TimeSpent
		r.TotalTimeSpent += cycle.TimeSpent

		if r.ExecutaionCap != 0 && len(r.Cycles) == r.ExecutaionCap {
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
