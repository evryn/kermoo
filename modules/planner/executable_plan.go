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
	templateValue common.SingleValueF
	targetValue   *float32
}

func (v *ExecutableValue) GetValue() float32 {
	if v.targetValue == nil {
		value, _ := v.templateValue.GetValue()
		v.targetValue = &value
	}

	return *v.targetValue
}

func NewExecutableValue(value common.SingleValueF) ExecutableValue {
	return ExecutableValue{
		templateValue: value,
	}
}
