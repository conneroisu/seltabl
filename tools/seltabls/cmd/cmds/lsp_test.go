package cmds

import (
	"testing"
)

func TestIsNull(t *testing.T) {
	var ptr *int
	var ch chan int
	var fn func()
	var mp map[string]int
	var sl []int
	var iface interface{}
	var strct struct{}

	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"nil pointer", ptr, true},
		{"nil channel", ch, true},
		{"nil function", fn, true},
		{"nil map", mp, true},
		{"nil slice", sl, true},
		{"nil interface", iface, true},
		{"non-nil struct", strct, false},
		{"non-nil int", 42, false},
		{"non-nil string", "hello", false},
		{"non-nil slice", []int{1, 2, 3}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNull(tt.input)
			if result != tt.expected {
				t.Errorf("isNull(%v) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
