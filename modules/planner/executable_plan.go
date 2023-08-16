package planner

import (
	"kermoo/modules/common"
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
	templateValue common.SingleFloat
	templateSize  common.SingleSize
	targetSize    *common.Size
	targetValue   *float32
}

func (v *ExecutableValue) GetValue() float32 {
	if v.targetValue == nil {
		value, _ := v.templateValue.ToFloat()
		v.targetValue = &value
	}

	return *v.targetValue
}

func (v *ExecutableValue) GetSize() common.Size {
	if v.targetSize == nil {
		value, _ := v.templateSize.ToSize()
		v.targetSize = &value
	}

	return *v.targetSize
}

func NewExecutableValue(value common.SingleFloat, size common.SingleSize) ExecutableValue {
	return ExecutableValue{
		templateValue: value,
		templateSize:  size,
	}
}
