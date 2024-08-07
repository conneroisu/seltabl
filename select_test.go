package seltabl

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

// TestSelect tests the Run function
func TestSelect(t *testing.T) {
	t.Run("Tables Test of Select", func(t *testing.T) {
		t.Parallel()
		args := []struct {
			name    string
			doc     string
			input   string
			wantErr bool
		}{
			{
				name: "Test Select",
				doc: `
				<table>
					<tr>
						<td>a</td>
						<td>b</td>
					</tr>
					<tr>
						<td>1</td>
						<td>2</td>
					</tr>
					<tr>
						<td>3</td>
						<td>4</td>
					</tr>
				</table>
			`,
				input: `
				<table>
					<tr>
						<td>a</td>
						<td>b</td>
					</tr>
					<tr>
						<td>1</td>
						<td>2</td>
					</tr>
					<tr>
						<td>3</td>
						<td>4</td>
					</tr>
				</table>
			`,
				wantErr: false,
			},
			{
				name: "Test Select with no data",
				doc: `
			`,
				input: `
			`,
				wantErr: true,
			},
		}
		for _, tt := range args {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				doc, err := goquery.NewDocumentFromReader(
					strings.NewReader(tt.doc),
				)
				if err != nil {
					t.Fatalf("failed to create document: %v", err)
				}
				cellValue := doc.Find("td")
				s := selector{control: ctlInnerTextSelector}
				_, err = s.Select(cellValue)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"Select() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
			})
		}
	})
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
		s := selector{control: ctlInnerTextSelector}
		// Call the Run method
		cellText, err := s.Select(divs)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Check the expected result
		expected := "Hello, World!"
		if cellText != expected {
			t.Errorf(
				"Expected cellText to be %q, but got %q",
				expected,
				cellText,
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
		s := selector{control: ctlAttrSelector, query: "href"}
		// Call the Run method
		cellText, err := s.Select(links)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Check the expected result
		if cellText != url {
			t.Errorf("Expected cellText to be %q, but got %q", url, cellText)
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
		s := selector{
			control: ctlInnerTextSelector,
			query:   ctlInnerTextSelector,
		}
		// Call the Select method
		_, err = s.Select(divs)
		if err == nil {
			t.Errorf("Expected error: %v", err)
		}
	})

	t.Run("Parallel Table Test of Select", func(t *testing.T) {
		t.Parallel()
		args := []struct {
			name    string
			input   string
			wantErr bool
		}{
			{
				name: "Test Select",
				input: `
				<table>
					<tr>
						<td>a</td>
						<td>b</td>
					</tr>
					<tr>
						<td>1</td>
						<td>2</td>
					</tr>
					<tr>
						<td>3</td>
						<td>4</td>
					</tr>
				</table>
			`,
				wantErr: false,
			},
			{
				name: "Test Select with no data",
				input: `
			`,
				wantErr: true,
			},
		}
		for _, tt := range args {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				doc, err := goquery.NewDocumentFromReader(
					strings.NewReader(tt.input),
				)
				if err != nil {
					t.Fatalf("failed to create document: %v", err)
				}
				cellValue := doc.Find("td")
				s := selector{control: ctlInnerTextSelector}
				_, err = s.Select(cellValue)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"Select() error = %v, wantErr %v",
						err,
						tt.wantErr,
					)
					return
				}
			})
		}
	})

}

func TestSelects(t *testing.T) {
	// Test case: ctlInnerTextSelector with a valid element
	{
		html := `<div>Test Text</div>`
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
		s := selector{control: ctlInnerTextSelector}
		cellValue := doc.Find("div")

		text, err := s.Select(cellValue)
		assert.NoError(t, err)
		assert.Equal(t, "Test Text", text)
	}

	// Test case: ctlAttrSelector with an existing attribute
	{
		html := `<div id="123">Test Text</div>`
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
		s := selector{control: ctlAttrSelector, query: "id"}
		cellValue := doc.Find("div")

		text, err := s.Select(cellValue)
		assert.NoError(t, err)
		assert.Equal(t, "123", text)
	}

	// Test case: ctlAttrSelector with a non-existing attribute
	{
		html := `<div>Test Text</div>`
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
		s := selector{control: ctlAttrSelector, query: "id"}
		cellValue := doc.Find("div")

		text, err := s.Select(cellValue)
		assert.Error(t, err)
		assert.Equal(t, "", text)
	}

	// Test case: Unsupported control value
	{
		html := `<div>Test Text</div>`
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
		s := selector{control: "UnsupportedSelector"}
		cellValue := doc.Find("div")

		text, err := s.Select(cellValue)
		assert.Error(t, err)
		assert.Equal(t, "", text)
	}
}
