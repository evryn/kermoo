package fluent

import (
	"fmt"
	"math/rand"
	"time"
)

// ParsedValue represents parsed values of generic types (int64, float64, time.Duration).
// It can contain singular values, ranges, or arrays of values.
type ParsedValue[T int64 | float64 | time.Duration] struct {
	values    []T
	isBetween bool

	cachedValue *T
}

// GetValue retrieves the value according to its type (singular, range, array).
// For singular, it simply returns the value.
// For ranges, it returns a random value between the range.
// For arrays, it returns a random element from the array.
func (p *ParsedValue[T]) GetValue() T {
	if len(p.values) == 1 {
		return p.values[0]
	}

	if p.isBetween {
		return p.getBetween(p.values[0], p.values[1])
	}

	return p.getRandomElement(p.values)
}

// GetCachedValue retrieves the cached value if available.
// If not, it updates the cache with a fresh value and returns it.
func (p *ParsedValue[T]) GetCachedValue() T {
	if p.cachedValue == nil {
		return p.GetUpdatedCacheValue()
	}

	return *p.cachedValue
}

// GetUpdatedCacheValue fetches a fresh value according to the type (singular, range, array)
// and updates the cache with that value before returning it.
func (p *ParsedValue[T]) GetUpdatedCacheValue() T {
	value := p.GetValue()
	p.cachedValue = &value

	return *p.cachedValue
}

// GetValues simply returns the parsed values as an array.
func (p *ParsedValue[T]) GetValues() []T {
	if p.isBetween {
		return []T{p.getBetween(p.values[0], p.values[1])}
	}

	return p.values
}

// IsRanged determines whether the value is a ranged one
func (p *ParsedValue[T]) IsRanged() bool {
	return p.isBetween
}

// GetRange returns the range as min and max values
func (p *ParsedValue[T]) GetRange() (T, T, error) {
	if !p.isBetween {
		return 0, 0, fmt.Errorf("value is not ranged")
	}

	return p.values[0], p.values[1], nil
}

// getBetween is a helper method that given a range, returns a random value falling between
// the specified min and max values.
func (p *ParsedValue[T]) getBetween(min, max T) T {
	return T(float64(min) + rand.Float64()*(float64(max)-float64(min)))
}

// getRandomElement is a helper method that returns a random element from the provided array of values.
func (p *ParsedValue[T]) getRandomElement(values []T) T {
	index := rand.Intn(len(values))
	return values[index]
}

// newParsedValue initializes and returns a ParsedValue object with the provided values and a flag
// indicating if the values represent a range or not.
func newParsedValue[T int64 | float64 | time.Duration](values []T, isBetween bool) ParsedValue[T] {
	return ParsedValue[T]{
		isBetween: isBetween,
		values:    values,
	}
}
