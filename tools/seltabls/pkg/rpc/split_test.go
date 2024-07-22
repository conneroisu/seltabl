package rpc

import (
	"bytes"
	"testing"
)

// TestSplit tests the Split function.
func TestSplit(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectedAdv int
		expectedTok []byte
		expectErr   bool
	}{
		{
			name:        "Valid Split",
			data:        []byte("Content-Length: 13\r\n\r\nHello, world!"),
			expectedAdv: 35,
			expectedTok: []byte("Content-Length: 13\r\n\r\nHello, world!"),
			expectErr:   false,
		},
		{
			name:        "Invalid Content Length",
			data:        []byte("Content-Length: abc\r\n\r\nHello, world!"),
			expectedAdv: 0,
			expectedTok: nil,
			expectErr:   true,
		},
		{
			name:        "Insufficient Content Length",
			data:        []byte("Content-Length: 20\r\n\r\nHello, world!"),
			expectedAdv: 0,
			expectedTok: nil,
			expectErr:   false,
		},
		{
			name:        "Missing Header",
			data:        []byte("Hello, world!"),
			expectedAdv: 0,
			expectedTok: nil,
			expectErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adv, tok, err := Split(tt.data, false)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if adv != tt.expectedAdv {
				t.Errorf("expected advance: %d, got: %d", tt.expectedAdv, adv)
			}
			if !bytes.Equal(tok, tt.expectedTok) {
				t.Errorf("expected token: %s, got: %s", tt.expectedTok, tok)
			}
		})
	}
}
