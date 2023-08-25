package planner_test

import (
	"kermoo/modules/fluent"
	"kermoo/modules/logger"
	"kermoo/modules/planner"

	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func teardownSubTest(t *testing.T) {
	Recorder.Reset()
}

var (
	name          = "My Plan"
	acceptedError = float64(0.1)
)

func TestSimplePlanExecution(t *testing.T) {
	logger.MustInitLogger("fatal")

	t.Run("executes simple plan with static percentage", func(t *testing.T) {
		defer teardownSubTest(t)

		plan := planner.NewPlan(planner.Plan{
			Interval: fluent.NewMustFluentDuration("10ms"),
			Duration: fluent.NewMustFluentDuration("50ms"),
			Name:     &name,
		})

		plan.Percentage = fluent.NewMustFluentFloat("50")

		plan.Assign(&Recorder)

		require.NoError(t, plan.Validate())

		plan.Start()

		Recorder.AssertTotalTimeSpent(t, 50*time.Millisecond, acceptedError)

		Recorder.AssertCycleValues(t, []ExpectedCycleValue{
			{Percentage: fluent.NewMustFluentFloat("50")},
			{Percentage: fluent.NewMustFluentFloat("50")},
			{Percentage: fluent.NewMustFluentFloat("50")},
			{Percentage: fluent.NewMustFluentFloat("50")},
			{Percentage: fluent.NewMustFluentFloat("50")},
		})
	})

	t.Run("executes simple plan with minimum and maximum value", func(t *testing.T) {
		defer teardownSubTest(t)

		plan := planner.NewPlan(planner.Plan{
			Interval: fluent.NewMustFluentDuration("10ms"),
			Duration: fluent.NewMustFluentDuration("50ms"),
			Name:     &name,
		})

		plan.Percentage = fluent.NewMustFluentFloat("10 to 90")

		plan.Assign(&Recorder)

		require.NoError(t, plan.Validate())

		plan.Start()

		// Assert that it took around 50ms (with 2ms error)
		Recorder.AssertTotalTimeSpent(t, 50*time.Millisecond, acceptedError)

		Recorder.AssertCycleValues(t, []ExpectedCycleValue{
			{Percentage: fluent.NewMustFluentFloat("10 to 90")},
			{Percentage: fluent.NewMustFluentFloat("10 to 90")},
			{Percentage: fluent.NewMustFluentFloat("10 to 90")},
			{Percentage: fluent.NewMustFluentFloat("10 to 90")},
			{Percentage: fluent.NewMustFluentFloat("10 to 90")},
		})
	})

	t.Run("executes simple plan with chart bar", func(t *testing.T) {
		defer teardownSubTest(t)

		plan := planner.NewPlan(planner.Plan{
			Interval: fluent.NewMustFluentDuration("10ms"),
			Duration: fluent.NewMustFluentDuration("50ms"),
			Name:     &name,
		})

		plan.Percentage = fluent.NewMustFluentFloat("0, 30, 70")

		plan.Assign(&Recorder)

		require.NoError(t, plan.Validate())

		plan.Start()

		Recorder.AssertTotalTimeSpent(t, 50*time.Millisecond, acceptedError)

		Recorder.AssertCycleValues(t, []ExpectedCycleValue{
			{Percentage: fluent.NewMustFluentFloat("0")},
			{Percentage: fluent.NewMustFluentFloat("30")},
			{Percentage: fluent.NewMustFluentFloat("70")},
			{Percentage: fluent.NewMustFluentFloat("0")},
			{Percentage: fluent.NewMustFluentFloat("30")},
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
		plan := planner.NewPlan(planner.Plan{
			Name: &name,
			SubPlans: []planner.SubPlan{
				{
					Percentage: fluent.NewMustFluentFloat("50"),
					Interval:   fluent.NewMustFluentDuration("10ms"),
					Duration:   fluent.NewMustFluentDuration("50ms"),
				},
				{
					Percentage: fluent.NewMustFluentFloat("10 to 90"),
					Interval:   fluent.NewMustFluentDuration("30ms"),
					Duration:   fluent.NewMustFluentDuration("60ms"),
				},
				{
					Percentage: fluent.NewMustFluentFloat("20, 30, 40"),
					Interval:   fluent.NewMustFluentDuration("10ms"),
					Duration:   fluent.NewMustFluentDuration("50ms"),
				},
			},
		})

		plan.Assign(&Recorder)

		require.NoError(t, plan.Validate())

		plan.Start()

		Recorder.AssertTotalTimeSpent(t,
			(50*time.Millisecond)+(60*time.Millisecond)+(50*time.Millisecond),
			acceptedError,
		)

		Recorder.AssertCycleValues(t, []ExpectedCycleValue{
			{Percentage: fluent.NewMustFluentFloat("50")},
			{Percentage: fluent.NewMustFluentFloat("50")},
			{Percentage: fluent.NewMustFluentFloat("50")},
			{Percentage: fluent.NewMustFluentFloat("50")},
			{Percentage: fluent.NewMustFluentFloat("50")},

			{Percentage: fluent.NewMustFluentFloat("10 to 90")},
			{Percentage: fluent.NewMustFluentFloat("10 to 90")},

			{Percentage: fluent.NewMustFluentFloat("20")},
			{Percentage: fluent.NewMustFluentFloat("30")},
			{Percentage: fluent.NewMustFluentFloat("40")},
			{Percentage: fluent.NewMustFluentFloat("20")},
			{Percentage: fluent.NewMustFluentFloat("30")},
		})
	})

	t.Run("executes sub plan with inifinit duration", func(t *testing.T) {
		defer teardownSubTest(t)
		plan := planner.NewPlan(planner.Plan{
			Name: &name,
			SubPlans: []planner.SubPlan{
				{
					Percentage: fluent.NewMustFluentFloat("50"),
					Interval:   fluent.NewMustFluentDuration("10ms"),
					Duration:   fluent.NewMustFluentDuration("50ms"),
				},
				{
					Percentage: fluent.NewMustFluentFloat("10 to 90"),
					Interval:   fluent.NewMustFluentDuration("30ms"),
					Duration:   fluent.NewMustFluentDuration("60ms"),
				},
				{
					Percentage: fluent.NewMustFluentFloat("20, 30, 40"),
					Interval:   fluent.NewMustFluentDuration("10ms"),
				},
			},
		})

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

		Recorder.AssertCycleValues(t, []ExpectedCycleValue{
			{Percentage: fluent.NewMustFluentFloat("50")},
			{Percentage: fluent.NewMustFluentFloat("50")},
			{Percentage: fluent.NewMustFluentFloat("50")},
			{Percentage: fluent.NewMustFluentFloat("50")},
			{Percentage: fluent.NewMustFluentFloat("50")},

			{Percentage: fluent.NewMustFluentFloat("10 to 90")},
			{Percentage: fluent.NewMustFluentFloat("10 to 90")},

			{Percentage: fluent.NewMustFluentFloat("20")},
			{Percentage: fluent.NewMustFluentFloat("30")},
			{Percentage: fluent.NewMustFluentFloat("40")},

			{Percentage: fluent.NewMustFluentFloat("20")},
			{Percentage: fluent.NewMustFluentFloat("30")},
			{Percentage: fluent.NewMustFluentFloat("40")},

			{Percentage: fluent.NewMustFluentFloat("20")},
			{Percentage: fluent.NewMustFluentFloat("30")},
			{Percentage: fluent.NewMustFluentFloat("40")},

			{Percentage: fluent.NewMustFluentFloat("20")},
			{Percentage: fluent.NewMustFluentFloat("30")},
			{Percentage: fluent.NewMustFluentFloat("40")},

			{Percentage: fluent.NewMustFluentFloat("20")},
		})
	})
}
