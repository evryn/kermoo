package planner

import (
	"kermoo/modules/values"
	"time"
)

type ExecutablePlan struct {
	Values       []*ExecutableValue
	Interval     time.Duration
	CurrentTries uint64
	TotalTries   uint64
	IsForever    bool
}

type ExecutableValue struct {
	templateValue values.SingleFloat
	templateSize  values.SingleSize
	targetSize    *values.Size
	targetValue   *float32
}

func (v *ExecutableValue) GetValue() float32 {
	if v.targetValue == nil {
		value, _ := v.templateValue.ToFloat()
		v.targetValue = &value
	}

	return *v.targetValue
}

func (v *ExecutableValue) GetSize() values.Size {
	if v.targetSize == nil {
		value, _ := v.templateSize.ToSize()
		v.targetSize = &value
	}

	return *v.targetSize
}

func NewExecutableValue(value values.SingleFloat, size values.SingleSize) ExecutableValue {
	return ExecutableValue{
		templateValue: value,
		templateSize:  size,
	}
}
