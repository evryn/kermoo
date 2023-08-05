package common

import (
	"buggybox/modules/utils"
	"fmt"
)

type SingleValueF struct {
	Exactly *float32  `json:"exactly"`
	Between []float32 `json:"between"`
}

func (s *SingleValueF) GetValue() (float32, error) {
	if s.Exactly != nil {
		return *s.Exactly, nil
	}

	if len(s.Between) != 2 {
		return 0, fmt.Errorf("value of `between` needs to have exactly two element as range")
	}

	min := s.Between[0]
	max := s.Between[1]

	if min > max {
		t := min
		min = max
		max = t
	}

	return utils.RandomFloat(min, max), nil
}
