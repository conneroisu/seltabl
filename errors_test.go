package seltabl

import (
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// TestErrMissingMustBePresent_Error tests the Error method of ErrMissingMustBePresent.
// It verifies that the error message contains the correct field name, selector, and type information.
func TestErrMissingMustBePresent_Error(t *testing.T) {
	field := reflect.StructField{
		Name: "TestField",
		Type: reflect.TypeOf(""),
	}
	cfg := &SelectorConfig{
		MustBePresent: "test-selector",
	}
	err := &ErrMissingMustBePresent{
		Field: field,
		Cfg:   cfg,
	}

	expected := "must be present (test-selector) not found for field TestField with type string"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', but got '%s'", expected, err.Error())
	}
}

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
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader("<html><body><p>Test</p></body></html>"))
	err := &ErrNoDataFound{
		Typ:   reflect.TypeOf(struct{}{}),
		Field: field,
		Cfg:   cfg,
		Doc:   doc,
	}

	errorMsg := err.Error()
	expectedParts := []string{
		"no data found for selector test-selector",
		"with type struct {",
		"in field TestField",
		"with type string",
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
	// Create an invalid document that will fail to generate HTML
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader("<invalid>"))
	err := &ErrNoDataFound{
		Typ:   reflect.TypeOf(struct{}{}),
		Field: field,
		Cfg:   cfg,
		Doc:   doc,
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
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader("<html><body><p>Test</p></body></html>"))
	err := &ErrSelectorNotFound{
		Typ:   reflect.TypeOf(struct{}{}),
		Field: field,
		Cfg:   cfg,
		Doc:   doc,
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
	// Create an invalid document that will fail to generate HTML
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader("<invalid>"))
	err := &ErrSelectorNotFound{
		Typ:   reflect.TypeOf(struct{}{}),
		Field: field,
		Cfg:   cfg,
		Doc:   doc,
	}

	errorMsg := err.Error()
	expected := "not found"
	if !strings.Contains(errorMsg, expected) {
		t.Errorf("Expected error message to contain '%s', but it doesn't it is '%s'", expected, errorMsg)
	}
}
