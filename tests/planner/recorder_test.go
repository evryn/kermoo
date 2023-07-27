package planner_test

import (
	"buggybox/modules/planner"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Execution struct {
	TimeSpent       time.Duration
	ExecutablePlan  planner.ExecutablePlan
	ExecutableValue planner.ExecutableValue
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

func (r *PlanRecorder) RecordPreSleep(ep planner.ExecutablePlan, ev planner.ExecutableValue) planner.PlanSignal {
	r.Executions = append(r.Executions, Execution{
		ExecutablePlan:  ep,
		ExecutableValue: ev,
	})

	return planner.PLAN_SIGNAL_CONTINUE
}

func (r *PlanRecorder) RecordPostSleep(startedAt time.Time, timeSpent time.Duration) planner.PlanSignal {
	r.Executions[len(r.Executions)-1].TimeSpent = timeSpent
	r.TotalTimeSpent += timeSpent

	if r.ExecutaionCap != 0 && len(r.Executions) == r.ExecutaionCap {
		return planner.PLAN_SIGNAL_TERMINATE
	}

	return planner.PLAN_SIGNAL_CONTINUE
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
