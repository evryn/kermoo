package process_test

import (
	"kermoo/modules/common"
	"kermoo/modules/process"
	"kermoo/modules/utils"
	"testing"
	"time"
)

// Mock implementation of planner.PlannableTrait for testing purposes
type MockPlannableTrait struct{}

func (m MockPlannableTrait) Plan() error {
	return nil
}

// Mock implementation of common.SingleValueDur for testing purposes
type MockSingleValueDur struct{}

func (d MockSingleValueDur) GetValue() (time.Duration, error) {
	return 0, nil
}

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
				Delay: &common.SingleValueDur{
					Exactly: utils.NewP[common.Duration](common.Duration(time.Second)),
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
					After: common.SingleValueDur{
						Exactly: utils.NewP[common.Duration](common.Duration(time.Second)),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid with bad delay duration (between with single value)",
			process: process.Process{
				Delay: &common.SingleValueDur{
					Between: []common.Duration{common.Duration(time.Second)},
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
					After: common.SingleValueDur{
						Between: []common.Duration{common.Duration(time.Second)},
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
					After: common.SingleValueDur{},
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
