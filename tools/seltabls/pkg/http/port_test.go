package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHomeDir tests the HomeDir function by asserting that it returns
// the correct a non-empty string.
func TestHomeDir(t *testing.T) {
	assert.NotEqual(t, "", HomeDir())
}
