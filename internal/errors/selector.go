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
type SelectorNotFound[T any] struct {
	Struct    T
	FieldName string
	Selector  string
}

// genStructReprAndHighlight returns a string representation of the golang struct
// and highlights the selector that was used to find the header row which was
// not found.
//
// It is used by the HeaderNotFoundError struct.
func selectorStructHighlight[T any](structPtr T, selector string) (string, error) {
	val := reflect.ValueOf(structPtr)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("expected struct pointer, got %s", val.Kind())
	}
	val = val.Elem()
	structType := val.Type()
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Selector: %s\n", selector))
	result.WriteString(fmt.Sprintf("type struct %s {\n", structType.Name()))
	for i := 0; i < val.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := val.Field(i)
		skv, err := GenStructTagString(field, selector)
		if err != nil {
			return "", fmt.Errorf("failed to generate struct tag string: %w", err)
		}
		result.WriteString(fmt.Sprintf("\t%s %v %s\n", field.Name, fieldValue.Type(), *skv))
	}
	_, err := result.WriteString("}")
	if err != nil {
		return "", fmt.Errorf("failed to write struct tag: %w", err)
	}
	fmt.Println(result.String())
	return result.String(), nil
}

// GenStructTagString returns a string representation of the struct tag
// for a struct field.
//
// It is used by the HeaderNotFoundError struct.
func GenStructTagString(field reflect.StructField, highlightSelector string) (*string, error) {
	var result strings.Builder
	var err error
	result.WriteString("`")
	// split on '"' s and iterate over them
	tags := string(field.Tag)
	re := regexp.MustCompile(`(\w+):"([^"]*)"`)
	matches := re.FindAllStringSubmatch(tags, -1)
	for _, match := range matches {
		key := match[1]
		value := match[2]
		if strings.Contains(value, highlightSelector) {
			_, err = result.WriteString(fmt.Sprintf(" %s:%s", key, "==\""+value+"==\""))
			if err != nil {
				return nil, fmt.Errorf("failed to write struct tag: %w", err)
			}
		} else {
			_, err = result.WriteString(fmt.Sprintf(" %v:\"%v\"", match[1], match[2]))
			if err != nil {
				return nil, fmt.Errorf("failed to write struct tag: %w", err)
			}
		}
	}
	_, err = result.WriteString("`")
	if err != nil {
		return nil, fmt.Errorf("failed to write struct tag: %w", err)
	}
	res := result.String()
	return &res, nil
}

// Error implements the error interface
func (e *SelectorNotFound[T]) Error() string {
	structString, err := selectorStructHighlight(e.Struct, e.Selector)
	if err != nil {
		return fmt.Errorf("failed to generate struct string while reporting that a selector was not found: %w", err).Error()
	}
	return fmt.Sprintf(
		"selector %s not found for field %s\n%s",
		e.Selector,
		e.FieldName,
		structString,
	)
}
