package seltabl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// SetStructField sets a struct field to a value
// It uses generics to specify the type of the struct
// and the field name
func SetStructField[T any](
	structPtr *T,
	fieldName string,
	cellValue *goquery.Selection,
	selector string,
) error {
	v := reflect.ValueOf(structPtr).Elem()
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("no such field: %s in struct", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot set field: %s", fieldName)
	}
	fieldType := field.Type().Kind()
	cellText := cellValue.Text()
	exists := true
	switch selector {
	case innerTextSelector:
		cellText = strings.TrimSpace(cellText)
		if cellValue.Length() == 0 {
			return fmt.Errorf("failed to find selector: %s", selector)
		}
	case attrSelector:
		cellText, exists = cellValue.Attr(selector)
		if !exists {
			return fmt.Errorf("failed to find selector: %s", selector)
		}
	default:
		print("default")
	}
	switch fieldType {
	case reflect.String:
		field.SetString(cellText)
	case reflect.Int:
		in, err := strconv.Atoi(cellText)
		if err != nil {
			return fmt.Errorf("failed to parse int: %s", err)
		}
		field.SetInt(int64(in))
	case reflect.Int8:
		in, err := strconv.Atoi(cellText)
		if err != nil {
			return fmt.Errorf("failed to parse int: %s", err)
		}
		field.SetInt(int64(in))
	case reflect.Int16:
		in, err := strconv.Atoi(cellText)
		if err != nil {
			return fmt.Errorf("failed to parse int: %s", err)
		}
		field.SetInt(int64(in))
	case reflect.Int32:
		in, err := strconv.Atoi(cellText)
		if err != nil {
			return fmt.Errorf("failed to parse int: %s", err)
		}
		field.SetInt(int64(in))
	case reflect.Int64:
		in, err := strconv.ParseInt(cellText, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse int: %s", err)
		}
		field.SetInt(in)
	case reflect.Uint:
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %s", err)
		}
		field.SetUint(in)
	case reflect.Uint8:
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %s", err)
		}
		field.SetUint(in)
	case reflect.Uint16:
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %s", err)
		}
		field.SetUint(in)
	case reflect.Uint32:
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %s", err)
		}
		field.SetUint(in)
	case reflect.Uint64:
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %s", err)
		}
		field.SetUint(in)
	case reflect.Float32:
		in, err := strconv.ParseFloat(cellText, 32)
		if err != nil {
			return fmt.Errorf("failed to parse float: %s", err)
		}
		field.SetFloat(in)
	case reflect.Float64:
		in, err := strconv.ParseFloat(cellText, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float: %s", err)
		}
		field.SetFloat(in)
	default:
		return fmt.Errorf("unsupported type: %s", fieldType)
	}
	return nil
}
