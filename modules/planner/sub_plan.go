package planner

import (
	"buggybox/config"
	"buggybox/modules/common"
	"time"
)

type SubPlan struct {
	Value    *common.MixedValueF
	Interval *time.Duration
	Duration *time.Duration
}

func (s *SubPlan) ToExecutablePlan() (*ExecutablePlan, error) {
	interval := config.Default.Planner.Interval

	var executableValues []*ExecutableValue

	singleValues, err := s.Value.ToSingleValues()

	if err != nil {
		return nil, err
	}

	for _, v := range singleValues {
		ev := NewExecutableValue(v)
		executableValues = append(executableValues, &ev)
	}

	if s.Interval != nil {
		interval = *s.Interval
	}

	ep := ExecutablePlan{
		Values:       executableValues,
		Interval:     interval,
		CurrentTries: 0,
		TotalTries:   0,
		IsForever:    false,
	}

	if s.Duration == nil {
		ep.IsForever = true
	} else {
		dur := *s.Duration
		ep.TotalTries = uint64(dur.Nanoseconds() / interval.Nanoseconds())
	}

	return &ep, nil
}
