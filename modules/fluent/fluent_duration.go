package fluent

import (
	"encoding/json"
	"time"
)

// FluentDuration represents a parsed time.Duration value.
// type FluentDuration ParsedValue[time.Duration]

type FluentDuration struct {
	input string
	pv    *ParsedValue[time.Duration]
}

func (f FluentDuration) Get() time.Duration {
	return f.pv.GetValue()
}

func (f FluentDuration) GetArray() []time.Duration {
	return f.pv.GetValues()
}

func (f FluentDuration) GetCached() time.Duration {
	return f.pv.GetCachedValue()
}

func (f FluentDuration) GetUpdatedCache() time.Duration {
	return f.pv.GetUpdatedCacheValue()
}

func (f FluentDuration) GetParsedValue() *ParsedValue[time.Duration] {
	return f.pv
}

func (f *FluentDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.input)
}

func (ff *FluentDuration) UnmarshalJSON(data []byte) error {
	input, err := prepareUnmarshalString(data)

	if err != nil {
		return err
	}

	newFF, err := NewFluentDuration(input)
	if err != nil {
		return err
	}

	*ff = *newFF
	return nil
}

// NewFluentDuration initializes and returns a FluentDuration object by parsing a string input into durations.
// It is able to recognize exact values, ranges, or arrays.
func NewFluentDuration(input string) (*FluentDuration, error) {
	parser := newParser(input)

	parsedValues, err := parser.GetDuations()

	if err != nil {
		return nil, err
	}

	fluent := FluentDuration{
		input: input,
		pv:    parsedValues,
	}

	return &fluent, nil
}

func NewMustFluentDuration(input string) *FluentDuration {
	f, _ := NewFluentDuration(input)

	return f
}
