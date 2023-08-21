package cpu_test

import (
	"kermoo/modules/memory"
	"kermoo/modules/values"
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
			Size: &values.MultiSize{
				SingleSize: values.SingleSize{
					Between: []values.Size{100},
				},
			},
		}
		err := leak.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "plan validation")
	})
}

func TestGetUid(t *testing.T) {
	c := &memory.MemoryLeak{}
	assert.Equal(t, "cpu-manager", c.GetName())
}

func TestHasCustomPlan(t *testing.T) {
	t.Run("with custom plan", func(t *testing.T) {
		leak := memory.MemoryLeak{
			Size: &values.MultiSize{
				SingleSize: values.SingleSize{
					Between: []values.Size{100, 200},
				},
			},
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
	size := values.MultiSize{
		SingleSize: values.SingleSize{
			Between: []values.Size{100, 200},
		},
	}
	duration := values.Duration(0)
	interval := values.Duration(0)

	leak := memory.MemoryLeak{
		Size:     &size,
		Duration: &duration,
		Interval: &interval,
	}

	plan := leak.MakeInlinePlan()

	assert.Equal(t, size, plan.Size)
	assert.Equal(t, duration, plan.Duration)
	assert.Equal(t, interval, plan.Interval)
}

func TestMakeDefaultPlan(t *testing.T) {
	c := &memory.MemoryLeak{}
	assert.Nil(t, c.MakeDefaultPlan())
}

func TestStartAndStopLeaking(t *testing.T) {
	load := &memory.MemoryLeak{}

	assert.Len(t, load.GetLeakedData(), 0)

	load.StartLeaking(values.Size(100))

	assert.Len(t, load.GetLeakedData(), 100)

	load.StopLeaking()

	assert.Len(t, load.GetLeakedData(), 0)
}
