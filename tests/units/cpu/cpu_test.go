package cpu_test

import (
	"kermoo/modules/cpu"
	"kermoo/modules/planner"
	"testing"

	// Replace with the actual package path for common
	"github.com/stretchr/testify/assert"
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
			Utilize: cpu.CpuUtilize{
				PlanRefs: []string{"ref1", "ref2"},
			},
		}
		err := cpu.Validate()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "plan refs")
	})

	t.Run("should return error when plan validation fails", func(t *testing.T) {
		plan := planner.InitPlan(planner.Plan{})
		cpu := cpu.Cpu{
			Utilize: cpu.CpuUtilize{
				Plan: &plan,
			},
		}
		err := cpu.Validate()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "plan validation")
	})
}

func TestGetUid(t *testing.T) {
	c := &cpu.Cpu{}
	assert.Equal(t, "cpu-manager", c.GetUid())
}

func TestHasCustomPlan(t *testing.T) {
	t.Run("with custom plan", func(t *testing.T) {
		c := &cpu.Cpu{
			Utilize: cpu.CpuUtilize{
				Plan: &planner.Plan{},
			},
		}
		assert.True(t, c.HasCustomPlan())
	})

	t.Run("without custom plan", func(t *testing.T) {
		c := &cpu.Cpu{}
		assert.False(t, c.HasCustomPlan())
	})
}

func TestGetDesiredPlanNames(t *testing.T) {
	planRefs := []string{"plan1", "plan2"}
	c := cpu.Cpu{
		Utilize: cpu.CpuUtilize{
			PlanRefs: planRefs,
		},
	}
	assert.Equal(t, planRefs, c.GetDesiredPlanNames())
}

func TestMakeCustomPlan(t *testing.T) {
	plan := &planner.Plan{}
	c := &cpu.Cpu{
		Utilize: cpu.CpuUtilize{
			Plan: plan,
		},
	}
	assert.Equal(t, plan, c.MakeCustomPlan())
}

func TestMakeDefaultPlan(t *testing.T) {
	c := &cpu.Cpu{}
	assert.Nil(t, c.MakeDefaultPlan())
}
