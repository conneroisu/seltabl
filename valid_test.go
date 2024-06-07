package seltabl

import (
	"reflect"
	"testing"
)

// TestIsTypeSupported tests the isTypeSupported function
func TestIsTypeSupported(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		reflectType reflect.Type
		want        bool
	}{
		{
			name:        "Invalid type (struct)",
			reflectType: reflect.TypeOf(TestStruct{}),
			want:        false,
		},
		{
			name:        "Invalid type",
			reflectType: reflect.TypeOf(testing.T{}),
			want:        false,
		},

		{
			name:        "Invalid type InternalTest",
			reflectType: reflect.TypeOf(testing.InternalTest{}),
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			if got := isTypeSupported(tt.reflectType); got != tt.want {
				t.Errorf("isValidType() = %v, want %v", got, tt.want)
			}
		})
	}
}
