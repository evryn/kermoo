package e2e_test

import (
	"kermoo/modules/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryLeakEndToEnd(t *testing.T) {
	t.Run("works with a simple dedicated leak plan", func(t *testing.T) {
		e2e := NewE2E(t)

		startMem, _ := utils.GetMemoryUsage()

		e2e.Start(`
            memoryLeak:
              interval: 100ms
              size: 100Mi
		`, 3*time.Second)

		// Wait a few while
		time.Sleep(500 * time.Millisecond)

		expectedMem := startMem + (100 * 1024 * 1024)

		currentMem, err := utils.GetMemoryUsage()
		require.NoError(t, err)

		assert.GreaterOrEqual(t, uint64(float64(expectedMem)*1.2), currentMem)
		assert.LessOrEqual(t, uint64(float64(expectedMem)*0.9), currentMem)

		e2e.Wait()

		e2e.RequireTimedOut()
	})

	t.Run("works with a referenced leak plan", func(t *testing.T) {
		e2e := NewE2E(t)

		startMem, _ := utils.GetMemoryUsage()

		e2e.Start(`
            plans:
            - name: leak
              interval: 100ms
              size: 100Mi
            memoryLeak:
              planRefs:
              - leak
		`, 3*time.Second)

		// Wait a few while
		time.Sleep(500 * time.Millisecond)

		expectedMem := startMem + (100 * 1024 * 1024)

		currentMem, err := utils.GetMemoryUsage()
		require.NoError(t, err)

		assert.GreaterOrEqual(t, uint64(float64(expectedMem)*1.2), currentMem)
		assert.LessOrEqual(t, uint64(float64(expectedMem)*0.9), currentMem)

		e2e.Wait()

		e2e.RequireTimedOut()
	})
}
