package seltabl

import (
	"fmt"
	"reflect"
)

// ErrMissingMustBePresent is an error for when a must be present selector is not found
type ErrMissingMustBePresent struct {
	Field         reflect.StructField
	MustBePresent string
}

// Error implements the error interface
func (e *ErrMissingMustBePresent) Error() string {
	return fmt.Sprintf("must be present not found for field %s with type %s", &e.Field.Name, &e.Field.Type)
}
