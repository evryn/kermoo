package Planner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Execution struct {
	TimeSpent       time.Duration
	ExecutablePlan  ExecutablePlan
	ExecutableValue ExecutableValue
}

type PlanRecorder struct {
	TotalTimeSpent time.Duration
	Executions     []Execution
	ExecutaionCap  int
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

func (r *PlanRecorder) RecordPreSleep(ep ExecutablePlan, ev ExecutableValue) PlanSignal {
	r.Executions = append(r.Executions, Execution{
		ExecutablePlan:  ep,
		ExecutableValue: ev,
	})

	return PLAN_SIGNAL_CONTINUE
}

func (r *PlanRecorder) RecordPostSleep(startedAt time.Time, timeSpent time.Duration) PlanSignal {
	r.Executions[len(r.Executions)-1].TimeSpent = timeSpent
	r.TotalTimeSpent += timeSpent

	if r.ExecutaionCap != 0 && len(r.Executions) == r.ExecutaionCap {
		return PLAN_SIGNAL_TERMINATE
	}

	return PLAN_SIGNAL_CONTINUE
}

func (r *PlanRecorder) AssertExpectedValues(t *testing.T, expectedValues []ExpectedExecutionValue) {
	assert.Len(t, r.Executions, len(expectedValues))

	for i, ev := range expectedValues {
		v := r.Executions[i].ExecutableValue

		if ev.IsBetween {
			assert.GreaterOrEqual(t, v.GetValue(), ev.Minimum)
			assert.LessOrEqual(t, v.GetValue(), ev.Maximum)
			assert.Equal(t, v.GetValue(), v.GetValue(), "Ranged values should always get similar result which is determined by the first attempt to retrieve the value.")
		} else {
			assert.Equal(t, ev.Static, v.GetValue())
		}
	}
}

func (r *PlanRecorder) AssertTotalTimeSpent(t *testing.T, expectedDuration time.Duration, expectedError time.Duration) {
	assert.LessOrEqual(t, r.TotalTimeSpent-expectedDuration, expectedError)
}

var Recorder PlanRecorder

func teardownSubTest(t *testing.T) {
	Recorder.Reset()
}

var (
	name          = "My Plan"
	interval_10ms = 10 * time.Millisecond
	interval_30ms = 30 * time.Millisecond
	interval_1s   = 1 * time.Second
	duration_50ms = 50 * time.Millisecond
	duration_60ms = 60 * time.Millisecond
	duration_5s   = 5 * time.Second
	float_0_1     = float32(0.1)
	float_0_3     = float32(0.3)
	float_0_5     = float32(0.5)
	float_0_9     = float32(0.9)
	float_5_5     = float32(5.5)
)

func TestSimplePlanExecution(t *testing.T) {
	t.Run("executes simple plan with static value", func(t *testing.T) {
		defer teardownSubTest(t)

		plan := Plan{
			Value: &Value{
				Static: &float_0_5,
			},
			Interval: &interval_10ms,
			Duration: &duration_50ms,
			Name:     &name,
			internal: &PlanInternal{},
		}

		plan.Execute(Callbacks{
			PreSleep:  Recorder.RecordPreSleep,
			PostSleep: Recorder.RecordPostSleep,
		})

		Recorder.AssertTotalTimeSpent(t, 50*time.Millisecond, 2*time.Millisecond)

		Recorder.AssertExpectedValues(t, []ExpectedExecutionValue{
			{Static: 0.5},
			{Static: 0.5},
			{Static: 0.5},
			{Static: 0.5},
			{Static: 0.5},
		})
	})

	t.Run("executes simple plan with minimum and maximum value", func(t *testing.T) {
		defer teardownSubTest(t)

		plan := Plan{
			Value: &Value{
				Minimum: &float_0_1,
				Maximum: &float_0_9,
			},
			Interval: &interval_10ms,
			Duration: &duration_50ms,
			Name:     &name,
			internal: &PlanInternal{},
		}

		plan.Execute(Callbacks{
			PreSleep:  Recorder.RecordPreSleep,
			PostSleep: Recorder.RecordPostSleep,
		})

		// Assert that it took around 50ms (with 2ms error)
		Recorder.AssertTotalTimeSpent(t, 50*time.Millisecond, 2*time.Millisecond)

		Recorder.AssertExpectedValues(t, []ExpectedExecutionValue{
			{IsBetween: true, Minimum: 0.1, Maximum: 0.9},
			{IsBetween: true, Minimum: 0.1, Maximum: 0.9},
			{IsBetween: true, Minimum: 0.1, Maximum: 0.9},
			{IsBetween: true, Minimum: 0.1, Maximum: 0.9},
			{IsBetween: true, Minimum: 0.1, Maximum: 0.9},
		})
	})

	t.Run("executes simple plan with chart bar", func(t *testing.T) {
		defer teardownSubTest(t)

		plan := Plan{
			Value: &Value{
				Chart: &Chart{
					Bars: []float32{0, 0.3, 0.7},
				},
			},
			Interval: &interval_10ms,
			Duration: &duration_50ms,
			Name:     &name,
			internal: &PlanInternal{},
		}

		plan.Execute(Callbacks{
			PreSleep:  Recorder.RecordPreSleep,
			PostSleep: Recorder.RecordPostSleep,
		})

		Recorder.AssertTotalTimeSpent(t, 50*time.Millisecond, 2*time.Millisecond)

		Recorder.AssertExpectedValues(t, []ExpectedExecutionValue{
			{Static: 0},
			{Static: 0.3},
			{Static: 0.7},
			{Static: 0},
			{Static: 0.3},
		})
	})
}

func TestSubPlanExecution(t *testing.T) {
	t.Run("fails when value is set along with sub plans", func(t *testing.T) {
		t.Skip("TODO: Implement")
	})

	t.Run("fails when interval is set along with sub plans", func(t *testing.T) {
		t.Skip("TODO: Implement")
	})

	t.Run("fails when duration is set along with sub plans", func(t *testing.T) {
		t.Skip("TODO: Implement")
	})

	t.Run("executes sub plan with specific duration", func(t *testing.T) {
		defer teardownSubTest(t)
		plan := Plan{
			Name:     &name,
			internal: &PlanInternal{},
			SubPlans: []SubPlan{
				{
					Value: &Value{
						Static: &float_0_5,
					},
					Interval: &interval_10ms,
					Duration: &duration_50ms,
				},
				{
					Value: &Value{
						Minimum: &float_0_1,
						Maximum: &float_0_9,
					},
					Interval: &interval_30ms,
					Duration: &duration_60ms,
				},
				{
					Value: &Value{
						Chart: &Chart{
							Bars: []float32{0.2, 0.3, 0.4},
						},
					},
					Interval: &interval_10ms,
					Duration: &duration_50ms,
				},
			},
		}

		plan.Execute(Callbacks{
			PreSleep:  Recorder.RecordPreSleep,
			PostSleep: Recorder.RecordPostSleep,
		})

		Recorder.AssertTotalTimeSpent(t,
			(50*time.Millisecond)+(60*time.Millisecond)+(50*time.Millisecond),
			5*time.Millisecond,
		)

		Recorder.AssertExpectedValues(t, []ExpectedExecutionValue{
			{Static: 0.5},
			{Static: 0.5},
			{Static: 0.5},
			{Static: 0.5},
			{Static: 0.5},
			{IsBetween: true, Minimum: 0.1, Maximum: 0.9},
			{IsBetween: true, Minimum: 0.1, Maximum: 0.9},
			{Static: 0.2},
			{Static: 0.3},
			{Static: 0.4},
			{Static: 0.2},
			{Static: 0.3},
		})
	})

	t.Run("executes sub plan with inifinit duration", func(t *testing.T) {
		defer teardownSubTest(t)
		plan := Plan{
			Name:     &name,
			internal: &PlanInternal{},
			SubPlans: []SubPlan{
				{
					Value: &Value{
						Static: &float_0_5,
					},
					Interval: &interval_10ms,
					Duration: &duration_50ms,
				},
				{
					Value: &Value{
						Minimum: &float_0_1,
						Maximum: &float_0_9,
					},
					Interval: &interval_30ms,
					Duration: &duration_60ms,
				},
				{
					Value: &Value{
						Chart: &Chart{
							Bars: []float32{0.2, 0.3, 0.4},
						},
					},
					Interval: &interval_10ms,
				},
			},
		}

		// Limit the execution to 20 times. It's here to preven the real inifnit number of executions -
		// just enough to test.
		Recorder.ExecutaionCap = 20

		plan.Execute(Callbacks{
			PreSleep:  Recorder.RecordPreSleep,
			PostSleep: Recorder.RecordPostSleep,
		})

		// Static runs 5 times for 50ms
		// Between runs 2 times for 60ms
		// Chart runs unlimited time but since capped, ends with 13 times which equals 130ms
		Recorder.AssertTotalTimeSpent(t,
			(50*time.Millisecond)+(60*time.Millisecond)+(13*10*time.Millisecond),
			5*time.Millisecond,
		)

		Recorder.AssertExpectedValues(t, []ExpectedExecutionValue{
			{Static: 0.5},
			{Static: 0.5},
			{Static: 0.5},
			{Static: 0.5},
			{Static: 0.5},
			{IsBetween: true, Minimum: 0.1, Maximum: 0.9},
			{IsBetween: true, Minimum: 0.1, Maximum: 0.9},
			{Static: 0.2},
			{Static: 0.3},
			{Static: 0.4},
			{Static: 0.2},
			{Static: 0.3},
			{Static: 0.4},
			{Static: 0.2},
			{Static: 0.3},
			{Static: 0.4},
			{Static: 0.2},
			{Static: 0.3},
			{Static: 0.4},
			{Static: 0.2},
		})

	})
}
