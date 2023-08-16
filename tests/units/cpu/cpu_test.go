package cpu_test

import (
	"kermoo/modules/cpu"
	"kermoo/modules/planner"
	"kermoo/modules/values"
	"testing"

	// Replace with the actual package path for values
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	t.Run("should return error when no plan or plan refs is set", func(t *testing.T) {
		cpu := cpu.Cpu{}
		err := cpu.Validate()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "plan")
	})

	t.Run("should return error when plan refs contain more than one element", func(t *testing.T) {
		cpu := cpu.Cpu{
			Load: cpu.CpuLoader{
				PlanRefs: []string{"ref1", "ref2"},
			},
		}
		err := cpu.Validate()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "plan refs")
	})

	t.Run("should return error when plan validation fails", func(t *testing.T) {
		plan := planner.NewPlan(planner.Plan{})
		plan.Percentage = &values.MultiFloat{
			SingleFloat: values.SingleFloat{
				Between: []float32{0.1},
			},
		}
		cpu := cpu.Cpu{
			Load: cpu.CpuLoader{
				Plan: &plan,
			},
		}
		err := cpu.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "plan validation")
	})
}

func TestGetUid(t *testing.T) {
	c := &cpu.Cpu{}
	assert.Equal(t, "cpu-manager", c.GetName())
}

func TestHasCustomPlan(t *testing.T) {
	t.Run("with custom plan", func(t *testing.T) {
		c := &cpu.Cpu{
			Load: cpu.CpuLoader{
				Plan: &planner.Plan{},
			},
		}
		assert.True(t, c.HasInlinePlan())
	})

	t.Run("without custom plan", func(t *testing.T) {
		c := &cpu.Cpu{}
		assert.False(t, c.HasInlinePlan())
	})
}

func TestGetDesiredPlanNames(t *testing.T) {
	planRefs := []string{"plan1", "plan2"}
	c := cpu.Cpu{
		Load: cpu.CpuLoader{
			PlanRefs: planRefs,
		},
	}
	assert.Equal(t, planRefs, c.GetDesiredPlanNames())
}

func TestMakeCustomPlan(t *testing.T) {
	plan := &planner.Plan{}
	c := &cpu.Cpu{
		Load: cpu.CpuLoader{
			Plan: plan,
		},
	}
	assert.Equal(t, plan, c.MakeInlinePlan())
}

func TestMakeDefaultPlan(t *testing.T) {
	c := &cpu.Cpu{}
	assert.Nil(t, c.MakeDefaultPlan())
}
