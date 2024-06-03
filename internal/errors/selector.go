package errors

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// SelectorNotFound is returned when a header selector is not found
// for a struct field
//
// It shows a generated golang stuct and highlights the header selector that was
// used to find the header row which was not found.
type SelectorNotFound[T any] struct {
	Struct        T
	FieldName     string
	hederSelector string
}

//
// func (e *HeaderNotFoundError[T]) Error() string {
//         return fmt.Sprintf(

// genStructReprAndHighlight returns a string representation of the golang struct
// and highlights the selector that was used to find the header row which was
// not found.
//
// It is used by the HeaderNotFoundError struct.
func genStructReprAndHighlight(structPtr interface{}, selector string) (string, error) {
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
		result.WriteString(fmt.Sprintf("\t%s %v %s\n", field.Name, fieldValue.Type(), genStructTagString(field, selector)))
	}
	result.WriteString("}")
	fmt.Println(result.String())
	return result.String(), nil
}

// genStructTagString returns a string representation of the struct tag
// for a struct field.
//
// It is used by the HeaderNotFoundError struct.
func genStructTagString(field reflect.StructField, highlightSelector string) string {
	var result strings.Builder
	result.WriteString("`")
	// split on '"' s and iterate over them
	tags := string(field.Tag)
	re := regexp.MustCompile(`(\w+):"([^"]*)"`)
	matches := re.FindAllStringSubmatch(tags, -1)
	for _, match := range matches {
		key := match[1]
		value := match[2]
		if strings.Contains(value, highlightSelector) {
			result.WriteString(fmt.Sprintf(" %s:%s", key, "==\""+value+"\"=="))
		} else {
			result.WriteString(fmt.Sprintf(" %v:\"%v\"", match[1], match[2]))
		}
	}
	result.WriteString("`")
	return result.String()
}
