package process_test

import (
	"kermoo/modules/fluent"
	"kermoo/modules/process"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcess_Validate(t *testing.T) {
	tests := []struct {
		name    string
		process process.Process
		wantErr bool
	}{
		{
			name: "valid default process",
			process: process.Process{
				Delay: nil,
				Exit:  nil,
			},
			wantErr: false,
		},
		{
			name: "valid with delay",
			process: process.Process{
				Delay: fluent.NewMustFluentDuration("1s"),
				Exit:  nil,
			},
			wantErr: false,
		},
		{
			name: "valid with exit",
			process: process.Process{
				Delay: nil,
				Exit: &process.ProcessExit{
					After: *fluent.NewMustFluentDuration("1s"),
				},
			},
			wantErr: false,
		},
		// {
		// 	name: "invalid with bad delay duration (between with single value)",
		// 	process: process.Process{
		// 		Delay: fluent.NewMustFluentDuration("1s to "),
		// 	},
		// 	wantErr: true,
		// },
		// {
		// 	name: "invalid with bad exit duration (between with single value)",
		// 	process: process.Process{
		// 		Exit: &process.ProcessExit{
		// 			After: *fluent.NewMustFluentDuration("1s to "),
		// 			Code:  0,
		// 		},
		// 	},
		// 	wantErr: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.process.Validate()
			if tt.wantErr {
				require.Error(t, err)
			}
		})
	}
}

func TestProcess_MakeCustomPlan(t *testing.T) {
	process := process.Process{
		Delay: nil,
		Exit: &process.ProcessExit{
			After: *fluent.NewMustFluentDuration("1s"),
		},
	}

	plan := process.MakeInlinePlan()

	assert.Equal(t, "process-manager", *plan.Name)
	assert.Equal(t, time.Second, plan.Duration.Get())
	assert.Equal(t, time.Second, plan.Interval.Get())
	assert.Equal(t, float64(100), plan.Percentage.Get())
}
