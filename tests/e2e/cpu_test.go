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
              utilize:
                plan:
                  interval: 1000ms
                  value:
                    exactly: 0.3
		`, 3*time.Second)

		// Wait a few while
		time.Sleep(500 * time.Millisecond)

		percentage, err := utils.GetCpuUsage(100 * time.Millisecond)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, float32(0.6), percentage)
		assert.LessOrEqual(t, float32(0.4), percentage)

		e2e.Wait()

		e2e.RequireTimedOut()
	})

}
