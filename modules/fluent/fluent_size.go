package fluent

import (
	"encoding/json"
)

// FluentSize is a human-friendly representation of a digital size (like memory size).
// You can specify them using: 50 (bytes), 100Ki (kilibytes), 100K (kilobytes), 5Mi, 60M, 1G, 1Gi, etc.
//
// Here are a few examples on how to define your desired value:
//
// - Define an specific size like "100Mi". I'll use it exactly as is.
//
// - Define a ranged size like "100Mi to 1Gi". I'll find a value randomly between them.
//
// - Define an array of size like "20Mi, 150K, 100, 1G". Some modules will pick one among them
// randomly or iterate over them like a graph of bars.
type FluentSize struct {
	input string
	pv    *ParsedValue[int64]
}

func (f FluentSize) Get() int64 {
	return f.pv.GetValue()
}

func (f FluentSize) GetArray() []int64 {
	return f.pv.GetValues()
}

func (f FluentSize) GetCached() int64 {
	return f.pv.GetCachedValue()
}

func (f FluentSize) GetUpdatedCache() int64 {
	return f.pv.GetUpdatedCacheValue()
}

func (f FluentSize) GetParsedValue() *ParsedValue[int64] {
	return f.pv
}

func (f *FluentSize) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.input)
}

func (ff *FluentSize) UnmarshalJSON(data []byte) error {
	input, err := prepareUnmarshalString(data)

	if err != nil {
		return err
	}

	newFF, err := NewFluentSize(input)
	if err != nil {
		return err
	}

	*ff = *newFF
	return nil
}

// NewFluentSize initializes and returns a FluentSize object by parsing a string input into floating point numbers.
// It is able to recognize exact values, ranges, or arrays.
func NewFluentSize(input string) (*FluentSize, error) {
	parser := newParser(input)

	parsedValues, err := parser.GetSizes()

	if err != nil {
		return nil, err
	}

	fluent := FluentSize{
		input: input,
		pv:    parsedValues,
	}

	return &fluent, nil
}

func NewMustFluentSize(input string) *FluentSize {
	f, _ := NewFluentSize(input)

	return f
}
