package values

import (
	"fmt"
	"kermoo/modules/utils"
)

type SingleSize struct {
	Exactly *Size  `json:"exactly"`
	Between []Size `json:"between"`
}

func (s *SingleSize) ToSize() (Size, error) {
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

	return Size(utils.RandomIntBetween(min.ToBytes(), max.ToBytes())), nil
}

func NewZeroSize() SingleSize {
	return SingleSize{
		Exactly: utils.NewP[Size](0),
	}
}
