package utils_test

import (
	"kermoo/modules/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYamlToJSON(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "Test valid yaml",
			input:    "name: John\nage: 30\n",
			expected: `{"age":30,"name":"John"}`,
			hasError: false,
		},
		{
			name:     "Test invalid yaml",
			input:    "\n: :\n",
			expected: "",
			hasError: true,
		},
		{
			name:     "Test empty yaml",
			input:    "",
			expected: `{}`,
			hasError: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := utils.YamlToJSON(testCase.input)

			if testCase.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, result)
			}
		})
	}
}
