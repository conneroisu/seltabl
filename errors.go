package seltabl

import (
	"fmt"
	"reflect"

	"github.com/PuerkitoBio/goquery"
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
	Doc   *goquery.Document
}

// Error implements the error interface
func (e *ErrNoDataFound) Error() string {
	doc, err := e.Doc.Html()
	if err != nil {
		return fmt.Sprintf("failed to get data rows html: %s", err)
	}
	return fmt.Sprintf(
		"no data found for selector %s with type %s in field %s with type %s\n html: %s",
		e.Cfg.QuerySelector,
		e.Typ,
		e.Field.Name,
		e.Field.Type,
		doc,
	)
}

// ErrSelectorNotFound is an error for when a selector is not found
type ErrSelectorNotFound struct {
	Typ   reflect.Type
	Field reflect.StructField
	Cfg   *SelectorConfig
}

// Error implements the error interface
func (e *ErrSelectorNotFound) Error() string {
	return fmt.Sprintf(
		"selector %s with type %s not found for field %s with type %s",
		e.Cfg.QuerySelector,
		e.Typ,
		e.Field.Name,
		e.Field.Type,
	)
}
