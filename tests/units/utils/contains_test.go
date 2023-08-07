package utils_test

import (
	"kermoo/modules/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	testCases := []struct {
		name     string
		arr      []string
		str      string
		expected bool
	}{
		{
			name:     "Exists",
			arr:      []string{"foo", "bar", "baz"},
			str:      "bar",
			expected: true,
		},
		{
			name:     "Does not exist",
			arr:      []string{"foo", "bar", "baz"},
			str:      "qux",
			expected: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, utils.Contains(tt.arr, tt.str))
		})
	}
}
