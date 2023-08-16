package process_test

import (
	"kermoo/modules/process"
	"kermoo/modules/utils"
	"kermoo/modules/values"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
				Delay: &values.SingleDuration{
					Exactly: utils.NewP[values.Duration](values.Duration(time.Second)),
				},
				Exit: nil,
			},
			wantErr: false,
		},
		{
			name: "valid with exit",
			process: process.Process{
				Delay: nil,
				Exit: &process.ProcessExit{
					After: values.SingleDuration{
						Exactly: utils.NewP[values.Duration](values.Duration(time.Second)),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid with bad delay duration (between with single value)",
			process: process.Process{
				Delay: &values.SingleDuration{
					Between: []values.Duration{values.Duration(time.Second)},
				},
				Exit: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid with bad exit duration (between with single value)",
			process: process.Process{
				Delay: nil,
				Exit: &process.ProcessExit{
					After: values.SingleDuration{
						Between: []values.Duration{values.Duration(time.Second)},
					},
					Code: 0,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid with exit code but no duration",
			process: process.Process{
				Delay: nil,
				Exit: &process.ProcessExit{
					After: values.SingleDuration{},
					Code:  0,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.process.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcess_MakeCustomPlan(t *testing.T) {
	process := process.Process{
		Delay: nil,
		Exit: &process.ProcessExit{
			After: values.SingleDuration{
				Exactly: utils.NewP[values.Duration](values.Duration(time.Second)),
			},
		},
	}

	plan := process.MakeInlinePlan()

	assert.Equal(t, "process-manager", *plan.Name)
	assert.Equal(t, time.Second, time.Duration(*plan.Duration))
	assert.Equal(t, time.Second, time.Duration(*plan.Interval))
	assert.Equal(t, float32(1.0), *plan.Value.Exactly)
}
