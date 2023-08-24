package cpu_test

import (
	"kermoo/modules/fluent"
	"kermoo/modules/memory"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	t.Run("should return error when no plan or plan refs is set", func(t *testing.T) {
		leak := memory.MemoryLeak{}
		err := leak.Validate()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "plan")
	})

	t.Run("should return error when plan refs contain more than one element", func(t *testing.T) {
		leak := memory.MemoryLeak{
			PlanRefs: []string{"ref1", "ref2"},
		}
		err := leak.Validate()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "plan refs")
	})

	t.Run("should return error when plan validation fails", func(t *testing.T) {
		leak := memory.MemoryLeak{
			Size: fluent.NewMustFluentSize("100 to "),
		}
		err := leak.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "leak specifications")
	})
}

func TestGetUid(t *testing.T) {
	c := &memory.MemoryLeak{}
	assert.Equal(t, "memory-leaker", c.GetName())
}

func TestHasCustomPlan(t *testing.T) {
	t.Run("with custom plan", func(t *testing.T) {
		leak := memory.MemoryLeak{
			Size: fluent.NewMustFluentSize("100 to 200"),
		}
		assert.True(t, leak.HasInlinePlan())
	})

	t.Run("without custom plan", func(t *testing.T) {
		c := &memory.MemoryLeak{}
		assert.False(t, c.HasInlinePlan())
	})
}

func TestGetDesiredPlanNames(t *testing.T) {
	planRefs := []string{"plan1", "plan2"}
	leak := memory.MemoryLeak{
		PlanRefs: planRefs,
	}
	assert.Equal(t, planRefs, leak.GetDesiredPlanNames())
}

func TestMakeCustomPlan(t *testing.T) {
	size := fluent.NewMustFluentSize("100 to 200")
	duration := fluent.NewMustFluentDuration("0")
	interval := fluent.NewMustFluentDuration("0")

	leak := memory.MemoryLeak{
		Size:     size,
		Duration: duration,
		Interval: interval,
	}

	plan := leak.MakeInlinePlan()

	assert.Equal(t, size, plan.Size)
	assert.Equal(t, int64(0), int64(plan.Duration.Get()))
	assert.Equal(t, int64(0), int64(plan.Interval.Get()))
}

func TestMakeDefaultPlan(t *testing.T) {
	c := &memory.MemoryLeak{}
	assert.Nil(t, c.MakeDefaultPlan())
}

func TestStartAndStopLeaking(t *testing.T) {
	load := &memory.MemoryLeak{}

	assert.Len(t, load.GetLeakedData(), 0)

	load.StartLeaking(100)

	assert.Len(t, load.GetLeakedData(), 100)

	load.StopLeaking()

	assert.Len(t, load.GetLeakedData(), 0)
}
