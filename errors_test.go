package seltabl

import (
	"reflect"
	"strings"
	"testing"
)

// TestErrNoDataFound_Error tests the Error method of ErrNoDataFound with a valid HTML document.
// It checks if the error message contains all the necessary information including the selector,
// field name, types, and the HTML content.
func TestErrNoDataFound_Error(t *testing.T) {
	field := reflect.StructField{
		Name: "TestField",
		Type: reflect.TypeOf(""),
	}
	cfg := &SelectorConfig{
		QuerySelector: "test-selector",
	}
	err := &ErrNoDataFound{
		Typ:   reflect.TypeOf(struct{}{}),
		Field: field,
		Cfg:   cfg,
	}

	errorMsg := err.Error()
	expectedParts := []string{
		"no data found",
		"struct",
	}

	for _, part := range expectedParts {
		if !strings.Contains(errorMsg, part) {
			t.Errorf("Expected error message to contain '%s', but it doesn't it is '%s'", part, errorMsg)
		}
	}
}

// TestErrNoDataFound_Error_FailedHtml tests the Error method of ErrNoDataFound when HTML generation fails.
// It verifies that the error message indicates a failure in getting the HTML content.
func TestErrNoDataFound_Error_FailedHtml(t *testing.T) {
	field := reflect.StructField{
		Name: "TestField",
		Type: reflect.TypeOf(""),
	}
	cfg := &SelectorConfig{
		QuerySelector: "test-selector",
	}
	err := &ErrNoDataFound{
		Typ:   reflect.TypeOf(struct{}{}),
		Field: field,
		Cfg:   cfg,
	}

	errorMsg := err.Error()
	expected := "no data found"
	if !strings.Contains(errorMsg, expected) {
		t.Errorf("Expected error message to contain '%s', but it doesn't it is '%s'", expected, errorMsg)
	}
}

// TestErrSelectorNotFound_Error tests the Error method of ErrSelectorNotFound with a valid HTML document.
// It ensures the error message includes the selector, field name, types, and HTML content.
func TestErrSelectorNotFound_Error(t *testing.T) {
	field := reflect.StructField{
		Name: "TestField",
		Type: reflect.TypeOf(""),
	}
	cfg := &SelectorConfig{
		QuerySelector: "test-selector",
	}
	err := &ErrSelectorNotFound{
		Typ:   reflect.TypeOf(struct{}{}),
		Field: field,
		Cfg:   cfg,
	}

	errorMsg := err.Error()
	expectedParts := []string{
		"selector test-selector",
		"with type struct {",
		"not found for field TestField",
		"with type string",
	}

	for _, part := range expectedParts {
		if !strings.Contains(errorMsg, part) {
			t.Errorf("Expected error message to contain '%s', but it doesn't", part)
		}
	}
}

// TestErrSelectorNotFound_Error_FailedHtml tests the Error method of ErrSelectorNotFound when HTML generation fails.
// It checks if the error message indicates a failure in getting the HTML content.
func TestErrSelectorNotFound_Error_FailedHtml(t *testing.T) {
	field := reflect.StructField{
		Name: "TestField",
		Type: reflect.TypeOf(""),
	}
	cfg := &SelectorConfig{
		QuerySelector: "test-selector",
	}
	err := &ErrSelectorNotFound{
		Typ:   reflect.TypeOf(struct{}{}),
		Field: field,
		Cfg:   cfg,
	}

	errorMsg := err.Error()
	expected := "not found"
	if !strings.Contains(errorMsg, expected) {
		t.Errorf("Expected error message to contain '%s', but it doesn't it is '%s'", expected, errorMsg)
	}
}

func TestErrNoDataFound_Error2(t *testing.T) {
	errNoDataFound := &ErrNoDataFound{
		Typ:   reflect.TypeOf("string"),
		Field: reflect.StructField{Name: "TestField", Type: reflect.TypeOf("string")},
		Cfg:   &SelectorConfig{QuerySelector: "div.invalid"},
	}

	expected := "failed to get data rows html: EOF"
	if strings.Contains(errNoDataFound.Error(), expected) {
		t.Errorf("ErrNoDataFound.Error() should not contain '%s', but it does", expected)
	}
}

// TestErrSelectorNotFound_Error2 tests the Error method of ErrSelectorNotFound
// with a valid HTML document.
//
// It ensures the error message includes the selector, field name, types, and
// HTML content.
func TestErrSelectorNotFound_Error2(t *testing.T) {
	errSelectorNotFound := &ErrSelectorNotFound{
		Typ:   reflect.TypeOf("string"),
		Field: reflect.StructField{Name: "TestField", Type: reflect.TypeOf("string")},
		Cfg:   &SelectorConfig{QuerySelector: "div.invalid"},
	}
	expected := "selector div.invalid with type string not found for field TestField with type string\n html: EOF"
	if strings.Contains(errSelectorNotFound.Error(), expected) {
		t.Errorf("ErrSelectorNotFound.Error() should not contain '%s', but it does", expected)
	}

}
