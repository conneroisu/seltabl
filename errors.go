package seltabl

import (
	"fmt"
	"reflect"
)

// ErrMissingMustBePresent is an error for when a must be present selector is not found
type ErrMissingMustBePresent struct {
	Field reflect.StructField
	Cfg   *SelectorConfig
}

// Error implements the error interface
func (e *ErrMissingMustBePresent) Error() string {
	return fmt.Sprintf(
		"must be present (%s) not found for field %s with type %s",
		e.Cfg.MustBePresent,
		e.Field.Name,
		e.Field.Type,
	)
}
