package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlwaysPasses(t *testing.T) {
	t.Log("This test always passes!")
	assert.True(t, true)
}
