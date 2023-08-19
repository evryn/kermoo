package utils_test

import (
	"kermoo/modules/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomInt64(t *testing.T) {
	tests := []struct {
		min, max int64
	}{
		{0, 10},
		{-10, 10},
		{-10, -1},
		{10, 10},
		{0, 1},
	}

	for _, tt := range tests {
		got := utils.RandomIntBetween(tt.min, tt.max)
		assert.GreaterOrEqual(t, tt.max, got)
		assert.LessOrEqual(t, tt.min, got)
	}
}
