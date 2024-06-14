package seltabl

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// SetStructField sets a struct field to a value
// It uses generics to specify the type of the struct
// and the field name.
//
// It is used by the NewFromString function.
func SetStructField[T any](
	structPtr *T,
	fieldName string,
	cellValue *goquery.Selection,
	selector SelectorInferface,
) error {
	var err error
	v := reflect.ValueOf(structPtr).Elem()
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("no such field: %s in struct", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot change the value of field: %s", fieldName)
	}
	fieldType := field.Type().Kind()
	// select the value from the cell
	value, err := selector.Select(cellValue)
	if err != nil {
		return fmt.Errorf("failed to run selector: %w", err)
	}
	// setting the field's value
	err = SetFieldValue(fieldType, value, &field)
	if err != nil {
		return fmt.Errorf("failed to insert value: %w", err)
	}
	return nil
}

// SetFieldValue sets the value of a field
//
// It is used by the SetStructField function to set the value of a struct field after
// selecting the value from a html node.
func SetFieldValue(
	fieldType reflect.Kind,
	cellText string,
	field *reflect.Value,
) error {
	switch fieldType {
	case reflect.String:
		field.SetString(cellText)
		return nil
	case reflect.Int:
		in, err := strconv.Atoi(cellText)
		if err != nil {
			return fmt.Errorf("failed to parse int: %s", err)
		}
		field.SetInt(int64(in))
		return nil
	case reflect.Int8:
		in, err := strconv.Atoi(cellText)
		if err != nil {
			return fmt.Errorf("failed to parse int: %s", err)
		}
		field.SetInt(int64(in))
		return nil
	case reflect.Int16:
		in, err := strconv.Atoi(cellText)
		if err != nil {
			return fmt.Errorf("failed to parse int: %s", err)
		}
		field.SetInt(int64(in))
		return nil
	case reflect.Int32:
		in, err := strconv.Atoi(cellText)
		if err != nil {
			return fmt.Errorf("failed to parse int: %s", err)
		}
		field.SetInt(int64(in))
		return nil
	case reflect.Int64:
		in, err := strconv.ParseInt(cellText, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse int: %s", err)
		}
		field.SetInt(in)
		return nil
	case reflect.Uint:
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %s", err)
		}
		field.SetUint(in)
		return nil
	case reflect.Uint8:
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %s", err)
		}
		field.SetUint(in)
		return nil
	case reflect.Uint16:
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %s", err)
		}
		field.SetUint(in)
		return nil
	case reflect.Uint32:
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %s", err)
		}
		field.SetUint(in)
		return nil
	case reflect.Uint64:
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %s", err)
		}
		field.SetUint(in)
		return nil
	case reflect.Float32:
		in, err := strconv.ParseFloat(cellText, 32)
		if err != nil {
			return fmt.Errorf("failed to parse float: %s", err)
		}
		field.SetFloat(in)
		return nil
	case reflect.Float64:
		in, err := strconv.ParseFloat(cellText, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float: %s", err)
		}
		field.SetFloat(in)
		return nil
	default:
		return fmt.Errorf("unsupported type: %s", fieldType)
	}
}
