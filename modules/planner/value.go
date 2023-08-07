package planner

import (
	"kermoo/config"
	"kermoo/modules/common"
)

type Value struct {
	Chart   *Chart
	Static  *float32
	Maximum *float32
	Minimum *float32
}

type Chart struct {
	Bars []float32
}

func (v Value) GetExecutableValues() []ExecutableValue {
	if v.Chart != nil && len(v.Chart.Bars) > 0 {
		var ev []ExecutableValue

		for _, bar := range v.Chart.Bars {
			ev = append(ev, NewExecutableValue(common.SingleValueF{
				Exactly: &bar,
			}))
		}

		return ev
	}

	min := config.Default.Planner.Minimum
	max := config.Default.Planner.Maximum

	if v.Maximum != nil {
		max = *v.Maximum
	}

	if v.Minimum != nil {
		min = *v.Minimum

		if min > max {
			min = max
		}
	}

	if v.Static != nil {
		min = *v.Static
		max = *v.Static
	}

	return []ExecutableValue{
		NewExecutableValue(common.SingleValueF{
			Between: []float32{min, max},
		}),
	}
}
