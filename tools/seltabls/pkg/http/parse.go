package http

import "fmt"

// ParseAddress parses a standard golang network address and returns the
// protocol and path.
//
// Example:
//
//	ip.ParseAddress("ipv4://127.0.0.1:80/hello")
//	// Output:
//	// Protocol: ipv4
//	// Path: 127.0.0.1:80/hello
func ParseAddress(addr string) (proto string, path string, err error) {
	m := netAddrRx.FindStringSubmatch(addr)
	if m == nil {
		return "", "", fmt.Errorf(
			"invalid address: %s", addr,
		)
	}
	return m[1], m[2], nil
}
