package errors

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// SelectionNotFound is returned when a header selector is not found
// for a struct field
//
// It shows a generated golang stuct and highlights the header selector that was
// used to find the header row which was not found.
type SelectionNotFound[T any] struct {
	Struct         T
	FieldName      string
	SelectionQuery string
}

// selectionStructHighlight returns a string representation of the golang struct
// and highlights the selector that was used to find the header row which was
// not found.
//
// It is used by the HeaderNotFoundError struct.
func selectionStructHighlight[T any](
	structPtr *T,
	selector string,
	fieldName string,
) (string, error) {
	val := reflect.ValueOf(structPtr)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("expected struct pointer, got %s", val.Kind())
	}
	val = val.Elem()
	structType := val.Type()
	var result strings.Builder
	_, err := result.WriteString(fmt.Sprintf("Selector: %s\n", selector))
	if err != nil {
		return "", fmt.Errorf("failed to write string: %w", err)
	}
	_, err = result.WriteString(fmt.Sprintf("type struct %s {\n", structType.Name()))
	if err != nil {
		return "", fmt.Errorf("failed to write string: %w", err)
	}
	for i := 0; i < val.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := val.Field(i)
		skv, err := genStructKeyString(field, selector, fieldName)
		if err != nil {
			return "", fmt.Errorf(
				"failed to generate struct key string: %w",
				err,
			)
		}
		_, err = result.WriteString(
			fmt.Sprintf("\t%s %v %s\n", field.Name, fieldValue.Type(), *skv),
		)
		if err != nil {
			return "", fmt.Errorf("failed to write string: %w", err)
		}
	}
	_, err = result.WriteString("}")
	if err != nil {
		return "", fmt.Errorf("failed to write string: %w", err)
	}
	return result.String(), nil
}

// genStructKeyString returns a string representation of the struct tag
// for a struct field.
//
// It is used by the SelectionNotFoundError struct.
func genStructKeyString(
	field reflect.StructField,
	highlightSelector string,
) (*string, error) {
	var result strings.Builder
	var err error
	current := false
	result.WriteString("`")
	// split on '"' s and iterate over them
	tags := string(field.Tag)
	re := regexp.MustCompile(`(\w+):"([^"]*)"`)
	matches := re.FindAllStringSubmatch(tags, -1)
	for _, match := range matches {
		key := match[1]
		value := match[2]
		if strings.Contains(value, highlightSelector) {
			if strings.Contains(key, field.Name) {
				current = true
			}
			_, err = result.WriteString(
				fmt.Sprintf(" %s:%s", key, "==\""+value+"\"=="),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to write string: %w", err)
			}
		} else {
			_, err = result.WriteString(fmt.Sprintf(" %v:\"%v\"", match[1], match[2]))
			if err != nil {
				return nil, fmt.Errorf("failed to write string: %w", err)
			}
		}
	}
	if !current {
		_, err = result.WriteString("`")
		if err != nil {
			return nil, fmt.Errorf("failed to write string: %w", err)
		}
	} else {
		_, err = result.WriteString("`     ❌")
		if err != nil {
			return nil, fmt.Errorf("failed to write string: %w", err)
		}
	}
	res := result.String()
	return &res, nil
}

// Error returns a string representation of the error
func (s *SelectionNotFound[T]) Error() string {
	val, err := selectionStructHighlight(
		&s.Struct,
		s.SelectionQuery,
		s.FieldName,
	)
	if err != nil {
		return fmt.Sprintf("failed to generate struct: %s", err)
	}
	return fmt.Sprintf(
		"no selection found for field %s with selector (%s)\n\n%s",
		s.FieldName,
		s.SelectionQuery,
		val,
	)
}
