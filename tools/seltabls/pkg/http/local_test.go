package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetLocalIP tests the GetLocalIP function by asserting that
// it returns a non-empty string.
//
// Prints the IP address to the buffer output of the test.
func TestGetLocalIP(t *testing.T) {
	t.Run("test get local ip", func(t *testing.T) {
		t.Parallel()
		ip, err := GetLocalIP()
		assert.NotEmpty(t, ip)
		assert.NoError(t, err)
		t.Logf("ip=%s", ip)
	})
}
