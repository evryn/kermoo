package planner

import (
	"kermoo/config"
	"kermoo/modules/common"
	"time"
)

type SubPlan struct {
	Value    *common.MixedValueF
	Interval *common.Duration
	Duration *common.Duration
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
		interval = time.Duration(*s.Interval)
	}

	ep := ExecutablePlan{
		Values:       executableValues,
		Interval:     time.Duration(interval),
		CurrentTries: 0,
		TotalTries:   0,
		IsForever:    false,
	}

	if s.Duration == nil {
		ep.IsForever = true
	} else {
		dur := time.Duration(*s.Duration)
		ep.TotalTries = uint64(dur.Nanoseconds() / interval.Nanoseconds())
	}

	return &ep, nil
}
