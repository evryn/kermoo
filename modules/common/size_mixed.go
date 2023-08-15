package common

import "fmt"

type MixedSize struct {
	SingleSize
	Chart *ChartSize `json:"chart"`
}

func (v MixedSize) ToSingleSizes() ([]SingleSize, error) {
	if v.Exactly != nil {
		return []SingleSize{{
			Exactly: v.Exactly,
		}}, nil
	}

	if v.Chart != nil && len(v.Chart.Bars) > 0 {
		var sv []SingleSize

		for i := range v.Chart.Bars {
			sv = append(sv, SingleSize{
				Exactly: &v.Chart.Bars[i],
			})
		}

		return sv, nil
	}

	if len(v.Between) == 2 {
		return []SingleSize{{
			Between: []Size{v.Between[0], v.Between[1]},
		}}, nil
	}

	if len(v.Between) != 2 {
		return nil, fmt.Errorf("value of `between` needs to have exactly two element as range")
	}

	return nil, fmt.Errorf("no value is set")
}
