package planner_test

import (
	"kermoo/modules/common"
	"kermoo/modules/logger"
	"kermoo/modules/planner"
	"kermoo/modules/utils"

	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func teardownSubTest(t *testing.T) {
	Recorder.Reset()
}

var (
	name          = "My Plan"
	interval_10ms = common.Duration(10 * time.Millisecond)
	interval_30ms = common.Duration(30 * time.Millisecond)
	duration_50ms = common.Duration(50 * time.Millisecond)
	duration_60ms = common.Duration(60 * time.Millisecond)
	float_0_1     = float32(0.1)
	float_0_5     = float32(0.5)
	float_0_9     = float32(0.9)
	acceptedError = float64(0.1)
)

func TestSimplePlanExecution(t *testing.T) {
	logger.MustInitLogger("fatal")

	t.Run("executes simple plan with static value", func(t *testing.T) {
		defer teardownSubTest(t)

		plan := planner.InitPlan(planner.Plan{
			Interval: &interval_10ms,
			Duration: &duration_50ms,
			Name:     &name,
		})

		plan.Value = &common.MultiFloat{
			SingleFloat: common.SingleFloat{
				Exactly: utils.NewP[float32](0.5),
			},
		}

		plan.Assign(&Recorder)

		require.NoError(t, plan.Validate())

		plan.Start()

		Recorder.AssertTotalTimeSpent(t, 50*time.Millisecond, acceptedError)

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

		plan := planner.InitPlan(planner.Plan{
			Interval: &interval_10ms,
			Duration: &duration_50ms,
			Name:     &name,
		})

		plan.Value = &common.MultiFloat{
			SingleFloat: common.SingleFloat{
				Between: []float32{0.1, 0.9},
			},
		}

		plan.Assign(&Recorder)

		require.NoError(t, plan.Validate())

		plan.Start()

		// Assert that it took around 50ms (with 2ms error)
		Recorder.AssertTotalTimeSpent(t, 50*time.Millisecond, acceptedError)

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

		plan := planner.InitPlan(planner.Plan{
			Interval: &interval_10ms,
			Duration: &duration_50ms,
			Name:     &name,
		})

		plan.Value = &common.MultiFloat{
			Chart: &common.FloatChart{
				Bars: []float32{0, 0.3, 0.7},
			},
		}

		plan.Assign(&Recorder)

		require.NoError(t, plan.Validate())

		plan.Start()

		Recorder.AssertTotalTimeSpent(t, 50*time.Millisecond, acceptedError)

		Recorder.AssertExpectedValues(t, []ExpectedExecutionValue{
			{Static: 0},
			{Static: 0.3},
			{Static: 0.7},
			{Static: 0},
			{Static: 0.3},
		})
	})

	t.Run("simple plan without duration lasts for ever", func(t *testing.T) {
		t.Skip("TODO: Implement")
	})
}

func TestSubPlanExecution(t *testing.T) {
	logger.MustInitLogger("fatal")

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
		plan := planner.InitPlan(planner.Plan{
			Name: &name,
			SubPlans: []planner.SubPlan{
				{
					Value:    &common.MultiFloat{},
					Interval: &interval_10ms,
					Duration: &duration_50ms,
				},
				{
					Value:    &common.MultiFloat{},
					Interval: &interval_30ms,
					Duration: &duration_60ms,
				},
				{
					Value:    &common.MultiFloat{},
					Interval: &interval_10ms,
					Duration: &duration_50ms,
				},
			},
		})

		plan.SubPlans[0].Value = &common.MultiFloat{
			SingleFloat: common.SingleFloat{
				Exactly: utils.NewP[float32](0.5),
			},
		}
		plan.SubPlans[1].Value = &common.MultiFloat{
			SingleFloat: common.SingleFloat{
				Between: []float32{0.1, 0.9},
			},
		}
		plan.SubPlans[2].Value = &common.MultiFloat{
			Chart: &common.FloatChart{
				Bars: []float32{0.2, 0.3, 0.4},
			},
		}

		plan.Assign(&Recorder)

		require.NoError(t, plan.Validate())

		plan.Start()

		Recorder.AssertTotalTimeSpent(t,
			(50*time.Millisecond)+(60*time.Millisecond)+(50*time.Millisecond),
			acceptedError,
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
		plan := planner.InitPlan(planner.Plan{
			Name: &name,
			SubPlans: []planner.SubPlan{
				{
					Value:    &common.MultiFloat{},
					Interval: &interval_10ms,
					Duration: &duration_50ms,
				},
				{
					Value:    &common.MultiFloat{},
					Interval: &interval_30ms,
					Duration: &duration_60ms,
				},
				{
					Value:    &common.MultiFloat{},
					Interval: &interval_10ms,
				},
			},
		})

		plan.SubPlans[0].Value = &common.MultiFloat{
			SingleFloat: common.SingleFloat{
				Exactly: utils.NewP[float32](0.5),
			},
		}
		plan.SubPlans[1].Value = &common.MultiFloat{
			SingleFloat: common.SingleFloat{
				Between: []float32{0.1, 0.9},
			},
		}
		plan.SubPlans[2].Value = &common.MultiFloat{
			Chart: &common.FloatChart{
				Bars: []float32{0.2, 0.3, 0.4},
			},
		}

		// Limit the execution to 20 times. It's here to preven the real inifnit number of executions -
		// just enough to test.
		Recorder.ExecutaionCap = 20

		plan.Assign(&Recorder)

		require.NoError(t, plan.Validate())

		plan.Start()

		// Static runs 5 times for 50ms
		// Between runs 2 times for 60ms
		// Chart runs unlimited time but since capped, ends with 13 times which equals 130ms
		Recorder.AssertTotalTimeSpent(t,
			(50*time.Millisecond)+(60*time.Millisecond)+(13*10*time.Millisecond),
			acceptedError,
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
