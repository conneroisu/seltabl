package struc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// increment increments a given integer
func increment(i int) int {
	return i + 1
}

// TestIncrement increments a given integer
func TestIncrement(t *testing.T) {
	a := assert.New(t)
	i := increment(1)
	a.Equal(2, i)
}
