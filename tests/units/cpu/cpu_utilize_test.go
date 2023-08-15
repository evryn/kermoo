package cpu_test

import (
	"kermoo/modules/cpu"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	cu := &cpu.CpuUtilize{}

	usage := float32(0.5)
	cu.Start(usage)

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
	cu := &cpu.CpuUtilize{}
	cu.Start(float32(0.5))

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
