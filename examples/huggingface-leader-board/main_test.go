package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRun tests the run function
func TestRun(t *testing.T) {
	err := run()
	assert.Nil(t, err)
}
