package errors

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
)

// SuperNovaStruct is a test struct
type SuperNovaStruct struct {
	Supernova string `json:"Supernova" seltabl:"Supernova" hSel:"tr:nth-child(1) th:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	Year      string `json:"Year"      seltabl:"Year"      hSel:"tr:nth-child(1) th:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
	Type      string `json:"Type"      seltabl:"Type"      hSel:"tr:nth-child(1) th:nth-child(3)" dSel:"tr td:nth-child(3)" cSel:"$text"`
	Distance  string `json:"Distance"  seltabl:"Distance"  hSel:"tr:nth-child(1) th:nth-child(4)" dSel:"tr td:nth-child(4)"`
	Notes     string `json:"Notes"     seltabl:"Notes"     hSel:"tr:nth-child(1) th:nth-child(5)" dSel:"tr td:nth-child(5)"`
}

// TestStructErrors tests the errors in the package
func TestStructErrors(t *testing.T) {
	stc := SuperNovaStruct{}
	output, err := selectionStructHighlight(
		&stc,
		"tr:nth-child(1) th:nth-child(1)",
	)
	assert.Nil(t, err)
	assert.NotEmpty(t, output)
	assert.True(t, strings.Contains(output, "Supernova"))
}

// Define a sample struct to be used in the tests
type SampleStruct struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

// TestSelectionStructHighlight tests the selectionStructHighlight function
func TestSelectionStructHighlight(t *testing.T) {
	sample := &SampleStruct{"John", 30, "john@example.com"}
	expected := "Selector: email\ntype struct SampleStruct {\n\tName string ` json:\"name\"`\n\tAge int ` json:\"age\"`\n\tEmail string ` json:==\"email\"==`\n}"
	result, err := selectionStructHighlight(sample, "email")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != expected {
		// diff the expected and result
		t.Logf("====================================================")
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(expected),
			B:        difflib.SplitLines(result),
			FromFile: "Expected",
			ToFile:   "Result",
			Context:  3,
		}
		d, err := difflib.GetUnifiedDiffString(diff)
		if err != nil {
			t.Fatalf("failed to generate diff: %v", err)
		}
		t.Errorf("diff:\n%s", d)
		t.Logf("====================================================")
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// TestGenStructKeyString tests the genStructKeyString function
func TestGenStructKeyString(t *testing.T) {
	field, _ := reflect.TypeOf(SampleStruct{}).FieldByName("Email")
	expected := "` json:==\"email\"==`"
	result, err := genStructKeyString(field, "email")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if *result != expected {
		t.Errorf("expected %v, got %v", expected, *result)
	}
}
