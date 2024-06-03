package errors

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/conneroisu/seltabl/internal"
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
func genStructReprAndHighlight(structPtr interface{}, selector string) string {
	val := reflect.ValueOf(structPtr)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return "Provided argument is not a pointer to a struct"
	}
	val = val.Elem()
	structType := val.Type()
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Selector: %s\n", selector))
	result.WriteString(fmt.Sprintf("type struct %s {\n", structType.Name()))
	for i := 0; i < val.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := val.Field(i).Interface()
		hSel := field.Tag.Get("hSel")
		rSel := internal.Warn(hSel)
		result.WriteString(fmt.Sprintf("\t%s: %v %s\n", field.Name, fieldValue, genStructTagString(field, rSel)))
	}
	result.WriteString("}")

	fmt.Println(result.String())
	print(result.String())

	return result.String()
}

// genStructTagString returns a string representation of the struct tag
// for a struct field.
//
// It is used by the HeaderNotFoundError struct.
func genStructTagString(field reflect.StructField, highlightSelector string) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf(" %s: %v `", field.Name, field.Type))
	// Split the tags by space and iterate over them
	tags := strings.Split(string(field.Tag), " ")
	for _, tag := range tags {
		if strings.HasPrefix(tag, highlightSelector) {
			internal.Warn(tag)
		}
		result.WriteString(fmt.Sprintf("%v ", tag))
	}

	result.WriteString("`")
	return result.String()
}
