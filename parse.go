package seltabl

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	innerTextSelector = "$text"
)

func NewFromString[T any](htmlInput string) ([]T, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlInput))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}
	dt := reflect.TypeOf((*T)(nil)).Elem()
	if dt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", dt.Kind())
	}
	results := new([]T)
	for i := 0; i < dt.NumField(); i++ {
		result := new(T)
		field := dt.Field(i)
		headName := field.Tag.Get("seltabl")
		if headName == "" {
			continue
		}
		headSelector := field.Tag.Get("hSel")
		if headSelector == "" {
			headSelector = innerTextSelector
		}
		headerRow := doc.Find(headSelector)
		if headerRow.Length() == 0 {
			return nil, fmt.Errorf("failed to find header row for field %s", headName)
		}
		dataSelector := field.Tag.Get("dSel")
		if dataSelector == "" {
			dataSelector = innerTextSelector
		}
		dataRow := doc.Find(dataSelector + ":not(" + headSelector + ")")
		if dataRow.Length() == 0 {
			return nil, fmt.Errorf("failed to find data row for field %s", headName)
		}
		fmt.Printf("headName: %s\n", headName)
		fmt.Printf("headSelector: %s\n", headSelector)
		fmt.Printf("headerRow: %s\n", headerRow.Text())
		fmt.Printf("dataSelector: %s\n", dataSelector)
		fmt.Printf("dataRow: %s\n", dataRow.Text())
		// surround each of the data cells with +++, so we can easily find them
		pieces := []string{}
		dataRow.Each(func(i int, s *goquery.Selection) {
			pieces = append(pieces, strings.TrimSpace(s.Text()))
			return
		})
		for i := 0; i < len(pieces); i++ {
			fmt.Printf("pieces[%d]: %s\n", i, pieces[i])
		}
		cell := headerRow.Eq(i)
		cellText := cell.Text()
		if err := SetStructField(result, field.Name, cellText); err != nil {
			return nil, fmt.Errorf("failed to set field %s: %s", field.Name, err)
		}
	}
	fmt.Printf("result: %v\n", *results)
	return *results, nil
}

func SetStructField[T any](structPtr *T, fieldName string, value interface{}) error {
	v := reflect.ValueOf(structPtr).Elem()

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("no such field: %s in struct", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot set field: %s", fieldName)
	}

	val := reflect.ValueOf(value)
	if field.Type() != val.Type() {
		return fmt.Errorf("provided value type didn't match struct field type")
	}
	field.Set(val)

	return nil
}
