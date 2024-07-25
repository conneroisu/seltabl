package seltabl

import (
	"fmt"
	"reflect"
)

// ErrNoDataFound is an error for when no data is found for a selector
type ErrNoDataFound struct {
	Typ   reflect.Type
	Field reflect.StructField
	Cfg   *SelectorConfig
}

// Error implements the error interface for ErrNoDataFound
func (e *ErrNoDataFound) Error() string {
	return fmt.Sprintf(
		"(%s) [%s] <%s> no data found for selector %s\n html",
		e.Typ,
		e.Field.Type,
		e.Field.Name,
		e.Cfg.QuerySelector,
	)
}

// ErrSelectorNotFound is an error for when a selector is not found
type ErrSelectorNotFound struct {
	Typ   reflect.Type        // type of the struct
	Field reflect.StructField // field of the struct
	Cfg   *SelectorConfig     // selector config
}

// Error implements the error interface for ErrSelectorNotFound
func (e *ErrSelectorNotFound) Error() string {
	return fmt.Sprintf(
		"selector %s with type %s not found for field %s with type %s\n html",
		e.Cfg.QuerySelector,
		e.Typ,
		e.Field.Name,
		e.Field.Type,
	)
}

// ErrParsing is returned when a field's value cannot be parsed.
type ErrParsing struct {
	Field reflect.Type
	Value string
	Err   error
}

// Error returns the error message. It implements the error interface.
func (e ErrParsing) Error() string {
	return fmt.Sprintf(
		"failed to parse field %s of type %s from value %s: %s",
		e.Field.Name(),
		e.Field.String(),
		e.Value,
		e.Err,
	)
}
