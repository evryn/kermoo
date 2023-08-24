package cpu_test

import (
	"kermoo/modules/cpu"
	"kermoo/modules/fluent"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	t.Run("should return error when no plan or plan refs is set", func(t *testing.T) {
		load := cpu.CpuLoader{}
		err := load.Validate()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "plan")
	})

	t.Run("should return error when plan refs contain more than one element", func(t *testing.T) {
		load := cpu.CpuLoader{
			PlanRefs: []string{"ref1", "ref2"},
		}
		err := load.Validate()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "plan refs")
	})

	t.Run("should return error when plan validation fails", func(t *testing.T) {
		load := cpu.CpuLoader{
			Percentage: fluent.NewMustFluentFloat("1 to "),
		}
		err := load.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "specification")
	})
}

func TestGetUid(t *testing.T) {
	c := &cpu.CpuLoader{}
	assert.Equal(t, "cpu-manager", c.GetName())
}

func TestHasCustomPlan(t *testing.T) {
	t.Run("with custom plan", func(t *testing.T) {
		load := cpu.CpuLoader{
			Percentage: fluent.NewMustFluentFloat("10 to 20"),
		}
		assert.True(t, load.HasInlinePlan())
	})

	t.Run("without custom plan", func(t *testing.T) {
		c := &cpu.CpuLoader{}
		assert.False(t, c.HasInlinePlan())
	})
}

func TestGetDesiredPlanNames(t *testing.T) {
	planRefs := []string{"plan1", "plan2"}
	load := cpu.CpuLoader{
		PlanRefs: planRefs,
	}
	assert.Equal(t, planRefs, load.GetDesiredPlanNames())
}

func TestMakeCustomPlan(t *testing.T) {
	percentage := fluent.NewMustFluentFloat("10 to 20")
	duration := fluent.NewMustFluentDuration("0")
	interval := fluent.NewMustFluentDuration("0")

	load := cpu.CpuLoader{
		Percentage: percentage,
		Duration:   duration,
		Interval:   interval,
	}

	plan := load.MakeInlinePlan()

	assert.Equal(t, percentage, plan.Percentage)
	assert.Equal(t, int64(0), int64(plan.Duration.Get()))
	assert.Equal(t, int64(0), int64(plan.Interval.Get()))
}

func TestMakeDefaultPlan(t *testing.T) {
	c := &cpu.CpuLoader{}
	assert.Nil(t, c.MakeDefaultPlan())
}

func TestStart(t *testing.T) {
	cu := &cpu.CpuLoader{}

	cu.Start(50)

	// Check that ctx and cancel are set
	ctx, cancel := cu.GetContextAndCancel()

	assert.NotNil(t, ctx)
	assert.NotNil(t, cancel)

	// TODO: Test CPU utilization

	// Check that the context has not been canceled
	select {
	case <-ctx.Done():
		t.Fatal("context should not be canceled")
	default:
	}
}

func TestStop(t *testing.T) {
	cu := &cpu.CpuLoader{}
	cu.Start(50)

	// Stop and check that the context is canceled
	cu.Stop()

	ctx, _ := cu.GetContextAndCancel()

	select {
	case <-ctx.Done():
		// Success, context is canceled
	default:
		t.Fatal("context should be canceled")
	}
}
