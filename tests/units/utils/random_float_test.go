package utils_test

import (
	"kermoo/modules/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomFloat(t *testing.T) {
	testCases := []struct {
		name string
		min  float32
		max  float32
	}{
		{
			name: "test with equal min and max",
			min:  5.0,
			max:  5.0,
		},
		{
			name: "test with non-zero min and max",
			min:  1.0,
			max:  10.0,
		},
		{
			name: "test with large range",
			min:  0.001,
			max:  100000.0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := utils.RandomFloatBetween(testCase.min, testCase.max)

			assert.GreaterOrEqual(t, testCase.max, result)
			assert.LessOrEqual(t, testCase.min, result)

			if testCase.max != testCase.min {
				result2 := utils.RandomFloatBetween(testCase.min, testCase.max)
				assert.NotEqual(t, result, result2, "random generator must return different value on each call")
			}
		})
	}
}
