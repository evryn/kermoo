package planner

import (
	"buggybox/modules/Utils"
	"time"
)

type ExecutablePlan struct {
	Values       []ExecutableValue
	Interval     time.Duration
	CurrentTries uint64
	TotalTries   uint64
	IsForever    bool
}

type ExecutableValue struct {
	Minimum     float32
	Maximum     float32
	targetValue *float32
}

func (v *ExecutableValue) GetValue() float32 {
	if v.targetValue == nil {
		if v.Maximum == v.Minimum {
			v.targetValue = &v.Maximum
		} else {
			rand := Utils.GenerateRandomFloat32Between(v.Minimum, v.Maximum)
			v.targetValue = &rand
		}
	}

	return *v.targetValue
}
