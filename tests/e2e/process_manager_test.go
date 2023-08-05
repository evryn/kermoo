package e2e_test

import (
	"testing"
	"time"
)

func TestProcessManagerEndToEnd(t *testing.T) {
	t.Run("exit with exact time", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            schemaVersion: "0.1-beta"
            process:
              exit:
                after:
                  exactly: 1s
                code: 20
		`, 5*time.Second)

		e2e.Wait()

		e2e.RequireNotTimedOut()
		e2e.AssertExitCode(20)
		e2e.AssertExecutaionDuration(1*time.Second, 2*time.Second)
	})

	t.Run("delayed exit with exact time", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            schemaVersion: "0.1-beta"
            process:
              exit:
                after:
                  exactly: 1s
                code: 20
              delay:
                exactly: 1s
		`, 5*time.Second)

		e2e.Wait()

		e2e.RequireNotTimedOut()
		e2e.AssertExitCode(20)
		e2e.AssertExecutaionDuration(2*time.Second, 3*time.Second)
	})

	t.Run("exit with ranged time", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            schemaVersion: "0.1-beta"
            process:
              exit:
                after:
                  between: ["1s", "2s"]
                code: 20
		`, 10*time.Second)

		e2e.Wait()

		e2e.RequireNotTimedOut()
		e2e.AssertExitCode(20)
		e2e.AssertExecutaionDuration(1*time.Second, 3*time.Second)
	})

	t.Run("test exit with delay", func(t *testing.T) {
		e2e := NewE2E(t)

		e2e.Start(`
            schemaVersion: "0.1-beta"
            process:
              exit:
                after:
                  between: ["1s", "2s"]
                code: 20
              delay:
                between: ["1s", "2s"]
		`, 10*time.Second)

		e2e.Wait()

		e2e.RequireNotTimedOut()
		e2e.AssertExitCode(20)
		e2e.AssertExecutaionDuration(2*time.Second, 5*time.Second)
	})
}
