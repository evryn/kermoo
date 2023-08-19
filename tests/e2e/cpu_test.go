package e2e_test

import (
	"kermoo/modules/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCpuEndToEnd(t *testing.T) {
	t.Run("works with a simple dedicated utilization plan", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            cpu:
              load:
                plan:
                  interval: 100ms
                  percentage:
                    exactly: 0.7
		`, 3*time.Second)

		// Wait a few while
		time.Sleep(500 * time.Millisecond)

		percentage, err := utils.GetCpuUsage(500 * time.Millisecond)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, float32(0.9), percentage)
		assert.LessOrEqual(t, float32(0.6), percentage)

		e2e.Wait()

		e2e.RequireTimedOut()
	})

	t.Run("works with a referenced utilization plan", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            plans:
            - name: spike
              interval: 100ms
              percentage:
                exactly: 0.7
            cpu:
              load:
                planRefs:
                - spike
		`, 3*time.Second)

		// Wait a few while
		time.Sleep(500 * time.Millisecond)

		percentage, err := utils.GetCpuUsage(500 * time.Millisecond)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, float32(0.9), percentage)
		assert.LessOrEqual(t, float32(0.6), percentage)

		e2e.Wait()

		e2e.RequireTimedOut()
	})

}
