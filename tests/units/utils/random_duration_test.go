package utils_test

import (
	"kermoo/modules/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomDuration(t *testing.T) {
	testCases := []struct {
		name     string
		min      time.Duration
		max      time.Duration
		hasError bool
	}{
		{
			name:     "Test with zero min and max",
			min:      0,
			max:      0,
			hasError: true,
		},
		{
			name: "test with non-zero min and max",
			min:  1 * time.Second,
			max:  10 * time.Second,
		},
		{
			name:     "test with min greater than max",
			min:      10 * time.Second,
			max:      1 * time.Second,
			hasError: true,
		},
		{
			name: "test with large range",
			min:  1 * time.Millisecond,
			max:  24 * time.Hour,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := utils.RandomDurationBetween(testCase.min, testCase.max)

			if testCase.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.GreaterOrEqual(t, testCase.max, *result)
				assert.LessOrEqual(t, testCase.min, *result)

				result2, _ := utils.RandomDurationBetween(testCase.min, testCase.max)
				assert.NotEqual(t, *result, *result2, "random value are supposed to be different on each call.")
			}
		})
	}
}
