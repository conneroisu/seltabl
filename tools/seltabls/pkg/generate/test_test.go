package generate

import (
	"bytes"
	"html/template"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
)

//go:embed test.tmpl
var testTmpl string

// TestTest generates a test from a given url.
func TestTestTemplateTest(t *testing.T) {
	a := assert.New(t)
	tmpl := template.New("test_file_template")
	tmpl, err := tmpl.Parse(testTmpl)
	if err != nil {
		t.Fatalf("Failed to parse test: %v", err)
	}
	// fill out the template
	args := struct {
		Name    string
		Version string
		URL     string
	}{
		Name:    "TestStruct",
		Version: "0.0.0.0",
		URL:     "https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "test", args)
	if err != nil {
		t.Fatalf("Failed to execute test template: %v", err)
	}
	// write the buffer to the file
	a.NotEqual(buf.String(), "")
	t.Logf("test: %s", buf.String())
}
