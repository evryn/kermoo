package values

import (
	"fmt"
	"kermoo/modules/utils"
)

type SingleFloat struct {
	Exactly *float32  `json:"exactly"`
	Between []float32 `json:"between"`
}

func (s *SingleFloat) ToFloat() (float32, error) {
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

	return utils.RandomFloatBetween(min, max), nil
}

func (s *SingleFloat) ToFloatRange() (float32, float32, error) {
	if s.Exactly != nil {
		return *s.Exactly, *s.Exactly, nil
	}

	if len(s.Between) != 2 {
		return 0, 0, fmt.Errorf("value of `between` needs to have exactly two element as range")
	}

	return s.Between[0], s.Between[1], nil
}

func NewZeroFloat() SingleFloat {
	return SingleFloat{
		Exactly: utils.NewP[float32](0),
	}
}
