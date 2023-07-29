package common

import (
	"buggybox/modules/utils"
	"fmt"
)

type SingleValueF struct {
	Exactly      *float32  `json:"exactly"`
	BetweenRange []float32 `json:"betweenRange"`
}

func (s *SingleValueF) GetValue() (float32, error) {
	if s.Exactly != nil {
		return *s.Exactly, nil
	}

	if len(s.BetweenRange) != 2 {
		return 0, fmt.Errorf("value of BetweenRange needs to have exactly two element as range")
	}

	min := s.BetweenRange[0]
	max := s.BetweenRange[1]

	if min > max {
		t := min
		min = max
		max = t
	}

	return utils.GenerateRandomFloat32Between(min, max), nil
}
