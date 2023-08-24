package fluent

import (
	"encoding/json"
)

// FluentFloat represents a parsed float64 value.
// type FluentFloat ParsedValue[float64]

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
