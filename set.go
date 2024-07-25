package seltabl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

// SetStructField sets a struct field to a value. It uses generics to specify
// the type of the struct and the field name. It also uses the selector
// interface to find the value and uses the type of the selector to parse and
// set the value.
//
// It is used by the NewFromString function.
func SetStructField[T any](
	structPtr *T,
	fieldName string,
	cellValue *goquery.Selection,
	selector SelectorI,
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
	err = setFieldValue(fieldType, value, &field)
	if err != nil {
		return fmt.Errorf("failed to insert value: %w", err)
	}
	return nil
}

// setFieldValue sets the value of a field
//
// It is used by the SetStructField function to set the value of a struct field
// after selecting the value from a html node.
//
// It ensures that the type of the field is compatible with the type of the
// value.
func setFieldValue(
	fieldType reflect.Kind,
	cellText string,
	field *reflect.Value,
) error {
	switch fieldType {
	case reflect.String:
		field.SetString(cellText)
		return nil
	case reflect.Int:
		if cellText == "" {
			cellText = "0"
		}
		in, err := strconv.Atoi(cellText)
		if err != nil {
			in, err = strconv.Atoi(extractNumbers(cellText))
			if err != nil {
				return ErrParsing{
					Field: field.Type(),
					Value: cellText,
					Err:   err,
				}
			}
		}
		field.SetInt(int64(in))
		return nil
	case reflect.Int8:
		if cellText == "" {
			cellText = "0"
		}
		in, err := strconv.Atoi(cellText)
		if err != nil {
			in, err = strconv.Atoi(extractNumbers(cellText))
			if err != nil {
				return ErrParsing{
					Field: field.Type(),
					Value: cellText,
					Err:   err,
				}
			}
		}
		field.SetInt(int64(in))
		return nil
	case reflect.Int16:
		if cellText == "" {
			cellText = "0"
		}
		in, err := strconv.Atoi(cellText)
		if err != nil {
			in, err = strconv.Atoi(extractNumbers(cellText))
			if err != nil {
				return ErrParsing{
					Field: field.Type(),
					Value: cellText,
					Err:   err,
				}
			}
		}
		field.SetInt(int64(in))
		return nil
	case reflect.Int32:
		if cellText == "" {
			cellText = "0"
		}
		in, err := strconv.Atoi(cellText)
		if err != nil {
			in, err = strconv.Atoi(extractNumbers(cellText))
			if err != nil {
				return ErrParsing{
					Field: field.Type(),
					Value: cellText,
					Err:   err,
				}
			}
		}
		field.SetInt(int64(in))
		return nil
	case reflect.Int64:
		if cellText == "" {
			cellText = "0"
		}
		in, err := strconv.ParseInt(cellText, 10, 64)
		if err != nil {
			in, err = strconv.ParseInt(extractNumbers(cellText), 10, 64)
			if err != nil {
				return ErrParsing{
					Field: field.Type(),
					Value: cellText,
					Err:   err,
				}
			}
		}
		field.SetInt(in)
		return nil
	case reflect.Uint:
		if cellText == "" {
			cellText = "0"
		}
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			in, err = strconv.ParseUint(extractNumbers(cellText), 10, 64)
			if err != nil {
				return ErrParsing{
					Field: field.Type(),
					Value: cellText,
					Err:   err,
				}
			}
		}
		field.SetUint(in)
		return nil
	case reflect.Uint8:
		if cellText == "" {
			cellText = "0"
		}
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			in, err = strconv.ParseUint(extractNumbers(cellText), 10, 64)
			if err != nil {
				return ErrParsing{
					Field: field.Type(),
					Value: cellText,
					Err:   err,
				}
			}
		}
		field.SetUint(in)
		return nil
	case reflect.Uint16:
		if cellText == "" {
			cellText = "0"
		}
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			in, err = strconv.ParseUint(extractNumbers(cellText), 10, 64)
			if err != nil {
				return ErrParsing{
					Field: field.Type(),
					Value: cellText,
					Err:   err,
				}
			}
		}
		field.SetUint(in)
		return nil
	case reflect.Uint32:
		if cellText == "" {
			cellText = "0"
		}
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			in, err = strconv.ParseUint(extractNumbers(cellText), 10, 64)
			if err != nil {
				return ErrParsing{
					Field: field.Type(),
					Value: cellText,
					Err:   err,
				}
			}
		}
		field.SetUint(in)
		return nil
	case reflect.Uint64:
		if cellText == "" {
			cellText = "0"
		}
		in, err := strconv.ParseUint(cellText, 10, 64)
		if err != nil {
			in, err = strconv.ParseUint(extractNumbers(cellText), 10, 64)
			if err != nil {
				return ErrParsing{
					Field: field.Type(),
					Value: cellText,
					Err:   err,
				}
			}
		}
		field.SetUint(in)
		return nil
	case reflect.Float32:
		if cellText == "" {
			cellText = "0"
		}
		in, err := strconv.ParseFloat(cellText, 32)
		if err != nil {
			in, err = strconv.ParseFloat(extractFloatNumbers(cellText), 32)
			if err != nil {
				return ErrParsing{
					Field: field.Type(),
					Value: cellText,
					Err:   err,
				}
			}
		}
		field.SetFloat(in)
		return nil
	case reflect.Float64:
		if cellText == "" {
			cellText = "0"
		}
		in, err := strconv.ParseFloat(cellText, 64)
		if err != nil {
			in, err = strconv.ParseFloat(extractFloatNumbers(cellText), 64)
			if err != nil {
				return ErrParsing{
					Field: field.Type(),
					Value: cellText,
					Err:   err,
				}
			}
		}
		field.SetFloat(in)
		return nil
	default:
		return fmt.Errorf("unsupported type: %s", fieldType)
	}
}

// reduceHTML removes all nodes from the selection that do not contain the
// text.
func reduceHTML(sel *goquery.Selection, text string) *goquery.Selection {
	var out *goquery.Selection
	for i := range sel.Length() {
		s := sel.Eq(i)
		body := s.Text()
		if !strings.Contains(body, text) {
			out = s.Remove()
		}
	}
	return out
}

// extractNumbers extracts all numbers from a string.
func extractNumbers(input string) string {
	var builder strings.Builder
	for _, char := range input {
		if unicode.IsDigit(char) {
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

// extractFloatNumbers extracts all numbers from a string including any decimal points.
func extractFloatNumbers(input string) string {
	var builder strings.Builder
	for _, char := range input {
		if unicode.IsDigit(char) || char == '.' {
			builder.WriteRune(char)
		}
	}
	return builder.String()
}
