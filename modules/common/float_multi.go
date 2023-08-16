package common

import "fmt"

type MultiFloat struct {
	SingleFloat
	Chart *FloatChart `json:"chart"`
}

func (v MultiFloat) ToSingleFloats() ([]SingleFloat, error) {
	if v.Exactly != nil {
		return []SingleFloat{{
			Exactly: v.Exactly,
		}}, nil
	}

	if v.Chart != nil && len(v.Chart.Bars) > 0 {
		var sv []SingleFloat

		for i := range v.Chart.Bars {
			sv = append(sv, SingleFloat{
				Exactly: &v.Chart.Bars[i],
			})
		}

		return sv, nil
	}

	if len(v.Between) == 2 {
		return []SingleFloat{{
			Between: []float32{v.Between[0], v.Between[1]},
		}}, nil
	}

	if len(v.Between) != 2 {
		return nil, fmt.Errorf("value of `between` needs to have exactly two element as range")
	}

	return nil, fmt.Errorf("no value is set")
}
