package seltabl

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// TestSelect tests the Run function
func TestSelect(t *testing.T) {
	t.Run("Select text from a div", func(t *testing.T) {
		// Create a mock cellValue
		cellValue, err := goquery.NewDocumentFromReader(
			strings.NewReader(
				"<html><body><div>Hello, World!</div></body></html>",
			),
		)
		if err != nil {
			t.Fatalf("failed to create document: %v", err)
		}
		divs := cellValue.Find("div")
		// Create a new instance of the selector
		s := selector{identifer: cSelInnerTextSelector}
		// Call the Run method
		cellText, err := s.Select(divs)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Check the expected result
		expected := "Hello, World!"
		if *cellText != expected {
			t.Errorf(
				"Expected cellText to be %q, but got %q",
				expected,
				*cellText,
			)
		}
	})

	t.Run("Select text from a href", func(t *testing.T) {
		// Create a mock cellValue
		url := "https://example.com"
		cellValue, err := goquery.NewDocumentFromReader(
			strings.NewReader(
				fmt.Sprintf(
					"<html><body><a href='%s'>Example</a></body></html>",
					url,
				),
			),
		)
		if err != nil {
			t.Fatalf("failed to create document: %v", err)
		}
		links := cellValue.Find("a")
		// Create a new instance of the selector
		s := selector{identifer: cSelAttrSelector, query: "href"}
		// Call the Run method
		cellText, err := s.Select(links)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Check the expected result
		if *cellText != url {
			t.Errorf("Expected cellText to be %q, but got %q", url, *cellText)
		}
	})

	t.Run("Select text from empty no text", func(t *testing.T) {
		t.Parallel()
		// Create a mock cellValue
		cellValue, err := goquery.NewDocumentFromReader(
			strings.NewReader("<html><body></body></html>"),
		)
		if err != nil {
			t.Fatalf("failed to create document: %v", err)
		}
		divs := cellValue.Find("div")
		// Create a new instance of the selector
		s := selector{identifer: cSelInnerTextSelector, query: cSelInnerTextSelector}
		// Call the Select method
		_, err = s.Select(divs)
		if err == nil {
			t.Errorf("Expected error: %v", err)
		}
	})

}
