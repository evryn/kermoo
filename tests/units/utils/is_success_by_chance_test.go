package utils_test

import (
	"buggybox/modules/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSuccessByChance(t *testing.T) {
	success := 0

	for i := 0; i < 100; i++ {
		if utils.IsSuccessByChance(0.5) == true {
			success = success + 1
		}
	}

	assert.GreaterOrEqual(t, 90, success)
	assert.LessOrEqual(t, 10, success)
}
