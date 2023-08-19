package planner_test

import (
	"kermoo/modules/planner"
	"kermoo/modules/values"
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

func (r *PlanRecorder) Reset() {
	r.Cycles = []planner.Cycle{}
	r.TotalTimeSpent = 0
}

func (r *PlanRecorder) AssertCycleValues(t *testing.T, expectedCycleValues []planner.CycleValue) {
	require.Len(t, r.Cycles, len(expectedCycleValues))

	for i, ev := range expectedCycleValues {
		actualPercentage, _ := r.Cycles[i].Value.Percentage.ToFloat()
		minPercentage, maxPercentage, _ := ev.Percentage.ToFloatRange()
		assert.GreaterOrEqual(t, maxPercentage, actualPercentage)
		assert.LessOrEqual(t, minPercentage, actualPercentage)

		actualSize, _ := r.Cycles[i].Value.Size.ToSize()
		minSize, maxSize, _ := ev.Size.ToSizeRange()
		assert.GreaterOrEqual(t, maxSize, actualSize)
		assert.LessOrEqual(t, minSize, actualSize)
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

func NewCycleValue(minPercentage, maxPercentage float32, minSize, maxSize values.Size) planner.CycleValue {
	return planner.CycleValue{
		Percentage: values.SingleFloat{
			Between: []float32{minPercentage, maxPercentage},
		},
		Size: values.SingleSize{
			Between: []values.Size{minSize, maxSize},
		},
	}
}

var Recorder PlanRecorder
