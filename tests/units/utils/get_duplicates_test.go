package utils_test

import (
	"buggybox/modules/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDuplicates(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "test with duplicates",
			input:    []string{"apple", "banana", "cherry", "apple", "banana"},
			expected: []string{"apple", "banana"},
		},
		{
			name:     "test without duplicates",
			input:    []string{"apple", "banana", "cherry"},
			expected: []string{},
		},
		{
			name:     "test with empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "test with all identical values",
			input:    []string{"apple", "apple", "apple", "apple"},
			expected: []string{"apple"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := utils.GetDuplicates(testCase.input)
			assert.Equal(t, testCase.expected, result)
			// if !reflect.DeepEqual(result, testCase.expected) {
			// 	t.Errorf("Expected %v, but got %v", testCase.expected, result)
			// }
		})
	}
}
