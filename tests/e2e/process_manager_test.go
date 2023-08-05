package e2e_test

import (
	"testing"
	"time"
)

func TestProcessManagerEndToEnd(t *testing.T) {

	t.Run("test exit and code", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            schemaVersion: "0.1-beta"
            process:
              exit:
                after:
                  exactly: 2s
                code: 20
		`, 5*time.Second)

		e2e.Wait()

		e2e.AssertExitCode(20)
		e2e.AssertExecutaionDuration(2*time.Second, 3*time.Second)
	})
}
