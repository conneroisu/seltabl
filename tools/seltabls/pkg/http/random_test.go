package http

import (
	"fmt"
	"net"
	"testing"
)

// TestRandomTCPPort tests the IsTCPPortAvailable function by
// asserting that it returns true when the port is available and
// false when the port is not available.
func TestRandomTCPPort(t *testing.T) {
	t.Run("Test Random TCP Port Available", func(t *testing.T) {
		p := RandomTCPPort()
		addr := fmt.Sprintf("127.0.0.1:%d", p)
		conn, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		conn.Close()
	})
	t.Run("Test Random TCP Port In Use", func(t *testing.T) {
		p := RandomTCPPort()
		addr := fmt.Sprintf("127.0.0.1:%d", p)
		conn, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		_, err2 := net.Listen("tcp", addr)
		if err2 == nil {
			t.Fatalf("addr should be in use %s", addr)
		}
	})
}
