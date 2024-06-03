package seltabl

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	innerTextSelector = "$text"
	attrSelector      = "@"
)

// NewFromString parses a string into a slice of structs.
// The struct must have a field with the tag `seltabl`, a header selector with the tag
// `hSel`, and a data selector with the tag `dSel`.
func NewFromString[T any](htmlInput string) ([]T, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlInput))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}
	dataType := reflect.TypeOf((*T)(nil)).Elem()
	if dataType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", dataType.Kind())
	}
	results := make([]T, 0)
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		headName := field.Tag.Get("seltabl")
		if headName == "" {
			continue
		}
		headSelector := field.Tag.Get("hSel")
		if headSelector == "" {
			return nil, fmt.Errorf(
				"failed to find header goquery selector for field %s",
				headName,
			)
		}
		headerRow := doc.Find(headSelector)
		if headerRow.Length() == 0 {
			return nil, fmt.Errorf(
				"failed to find header row for field %s",
				headName,
			)
		}
		dataSelector := field.Tag.Get("dSel")
		if dataSelector == "" {
			return nil, fmt.Errorf(
				"failed to find data goquery selector for field %s",
				headName,
			)
		}
		cellSelector := field.Tag.Get("cSel")
		if cellSelector == "" {
			cellSelector = innerTextSelector
		}
		dataRows := doc.Find(fmt.Sprintf("%s:not(%s)", dataSelector, headSelector))
		if dataRows.Length() == 0 {
			return nil, fmt.Errorf(
				"failed to find data row for field %s",
				headName,
			)
		}
		if len(results) == 0 {
			results = make([]T, dataRows.Length())
		}
		for j := 0; j < dataRows.Length(); j++ {
			if err := SetStructField(
				&results[j],
				field.Name,
				dataRows.Eq(j),
				cellSelector,
			); err != nil {
				return nil, fmt.Errorf(
					"failed to set field %s: %s",
					field.Name,
					err,
				)
			}
		}
	}
	return results, nil
}
