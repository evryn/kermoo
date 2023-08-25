package fluent

import (
	"encoding/json"
)

// FluentFloat is a human-friendly representation of a float amount like percentage.
// You can specify them like: 0, 2.5, 100, 60.5, ...
//
// Here are a few examples on how to define your desired value:
//
// - Define an specific duration like "20.8". I'll use it exactly as is.
//
// - Define a ranged duration like "5.2 to 40". I'll find a value randomly between them.
//
// - Define an array of durations like "1.5, 20, 50, 0". Some modules will pick one among them
// randomly or iterate over them like a graph of bars.
type FluentFloat struct {
	input string
	pv    *ParsedValue[float64]
}

func (f FluentFloat) Get() float64 {
	return f.pv.GetValue()
}

func (f FluentFloat) GetArray() []float64 {
	return f.pv.GetValues()
}

func (f FluentFloat) GetCached() float64 {
	return f.pv.GetCachedValue()
}

func (f FluentFloat) GetUpdatedCache() float64 {
	return f.pv.GetUpdatedCacheValue()
}

func (f FluentFloat) GetParsedValue() *ParsedValue[float64] {
	return f.pv
}

func (f *FluentFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.input)
}

func (ff *FluentFloat) UnmarshalJSON(data []byte) error {
	input, err := prepareUnmarshalString(data)

	if err != nil {
		return err
	}

	newFF, err := NewFluentFloat(input)
	if err != nil {
		return err
	}

	*ff = *newFF
	return nil
}

// NewFluentFloat initializes and returns a FluentFloat object by parsing a string input into floating point numbers.
// It is able to recognize exact values, ranges, or arrays.
func NewFluentFloat(input string) (*FluentFloat, error) {
	parser := newParser(input)

	parsedValues, err := parser.GetFloats()

	if err != nil {
		return nil, err
	}

	fluent := FluentFloat{
		input: input,
		pv:    parsedValues,
	}

	return &fluent, nil
}

func NewMustFluentFloat(input string) *FluentFloat {
	f, _ := NewFluentFloat(input)

	return f
}
