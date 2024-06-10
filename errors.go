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

// ErrNoDataFound is an error for when no data is found for a selector
type ErrNoDataFound struct {
	Typ   reflect.Type
	Field reflect.StructField
	Cfg   *SelectorConfig
}

// Error implements the error interface
func (e *ErrNoDataFound) Error() string {
	return fmt.Sprintf(
		"no data found for selector %s with type %s in field %s with type %s",
		e.Cfg.QuerySelector,
		e.Typ,
		e.Field.Name,
		e.Field.Type,
	)
}
