package common

import (
	"fmt"
	"kermoo/modules/utils"
	"time"
)

type SingleDuration struct {
	Exactly *Duration  `json:"exactly"`
	Between []Duration `json:"between"`
}

func (s *SingleDuration) ToStandardDuration() (time.Duration, error) {

	if s.Exactly != nil {
		return time.Duration(*s.Exactly), nil
	}

	if len(s.Between) != 2 {
		return 0, fmt.Errorf("value of `between` needs to have exactly two element as range")
	}

	min := time.Duration(s.Between[0])
	max := time.Duration(s.Between[1])

	if min > max {
		t := min
		min = max
		max = t
	}

	dur, err := utils.RandomDurationBetween(min, max)

	if err != nil {
		return 0, err
	}

	return *dur, nil
}
