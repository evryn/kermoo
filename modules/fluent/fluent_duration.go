package fluent

import (
	"encoding/json"
	"time"
)

// FluentDuration is a human-friendly representation of a duration. Use golang-format durations
// to define them - such as: 1h (hour), 2m (minute), 5s (second), 100ms (milisecond),
// 500ns (nanosecond), etc. You can also combine them too: 1m8s100ms
//
// Here are a few examples on how to define your desired duration:
//
// - Define an specific duration like "100ms". I'll use it exactly as is.
//
// - Define a ranged duration like "200ms to 2s". I'll find a value randomly between them.
//
// - Define an array of durations like "1s, 200ms, 3m100ms". I'll pick one randomly.
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
