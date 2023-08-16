package planner

import (
	"fmt"
	"kermoo/config"
	"kermoo/modules/common"
	"time"
)

type SubPlan struct {
	Value    *common.MultiFloat `json:"value"`
	Size     *common.MultiSize  `json:"size"`
	Interval *common.Duration   `json:"interval"`
	Duration *common.Duration   `json:"duration"`
}

func (s *SubPlan) BuildExecutableValues() ([]*ExecutableValue, error) {
	count := 0
	var err error
	var executableValues []*ExecutableValue
	var singleSizes []common.SingleSize
	var singleValues []common.SingleFloat

	if s.Size != nil {
		singleSizes, err = s.Size.ToSingleSizes()

		if err != nil {
			return nil, fmt.Errorf("failed to convert size to single values: %v", err)
		}

		if len(singleSizes) > count {
			count = len(singleSizes)
		}
	}

	if s.Value != nil {
		singleValues, err = s.Value.ToSingleValues()

		if err != nil {
			return nil, fmt.Errorf("failed to convert value to single values: %v", err)
		}

		if len(singleValues) > count {
			count = len(singleValues)
		}
	}

	if len(singleSizes) > 0 && len(singleValues) > 0 && len(singleSizes) != len(singleValues) {
		return nil, fmt.Errorf("both size and values are set while the count of individual steps does not match together")
	}

	for i := 0; i < count; i++ {
		value := common.MakeZeroSingleValueF()
		size := common.MakeZeroSize()

		if len(singleValues) >= i+1 {
			value = singleValues[i]
		}

		if len(singleSizes) >= i+1 {
			size = singleSizes[i]
		}

		ev := NewExecutableValue(value, size)
		executableValues = append(executableValues, &ev)
	}

	return executableValues, nil
}

func (s *SubPlan) ToExecutablePlan() (*ExecutablePlan, error) {
	interval := config.Default.Planner.Interval

	executableValues, err := s.BuildExecutableValues()

	if err != nil {
		return nil, fmt.Errorf("error building executable values: %v", err)
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
