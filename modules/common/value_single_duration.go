package common

import (
	"buggybox/modules/utils"
	"fmt"
	"time"
)

type SingleValueDur struct {
	Exactly      *Duration  `json:"exactly"`
	BetweenRange []Duration `json:"betweenRange"`
}

func (s *SingleValueDur) GetValue() (time.Duration, error) {

	if s.Exactly != nil {
		return time.Duration(*s.Exactly), nil
	}

	if len(s.BetweenRange) != 2 {
		return 0, fmt.Errorf("value of BetweenRange needs to have exactly two element as range")
	}

	min := time.Duration(s.BetweenRange[0])
	max := time.Duration(s.BetweenRange[1])

	if min > max {
		t := min
		min = max
		max = t
	}

	dur, err := utils.RandomDuration(min, max)

	if err != nil {
		return 0, err
	}

	return *dur, nil
}
