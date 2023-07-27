package planner

import (
	"buggybox/config"
	"time"
)

type SubPlan struct {
	Value    *Value
	Interval *time.Duration
	Duration *time.Duration
}

func (s *SubPlan) ToExecutablePlan() ExecutablePlan {
	interval := config.Default.Planner.Interval

	if s.Interval != nil {
		interval = *s.Interval
	}

	ep := ExecutablePlan{
		Values:       s.Value.GetExecutableValues(),
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

	return ep
}
