package common

import "fmt"

type MixedValueF struct {
	SingleValueF
	Chart *Chart `json:"chart"`
}

func (v MixedValueF) ToSingleValues() ([]SingleValueF, error) {
	if v.Exactly != nil {
		return []SingleValueF{{
			Exactly: v.Exactly,
		}}, nil
	}

	if v.Chart != nil && len(v.Chart.Bars) > 0 {
		var sv []SingleValueF

		for i := range v.Chart.Bars {
			sv = append(sv, SingleValueF{
				Exactly: &v.Chart.Bars[i],
			})
		}

		return sv, nil
	}

	if len(v.Between) == 2 {
		return []SingleValueF{{
			Between: []float32{v.Between[0], v.Between[1]},
		}}, nil
	}

	if len(v.Between) != 2 {
		return nil, fmt.Errorf("value of `between` needs to have exactly two element as range")
	}

	return nil, fmt.Errorf("no value is set")
}
