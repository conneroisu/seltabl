package generate

import (
	"bytes"
	"html/template"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
)

//go:embed struct.tmpl
var structTmpl string

// TestStruct generates a struct from a given url
func TestStructTemplateStruct(t *testing.T) {
	a := assert.New(t)
	tmpl := template.New("struct_file_template")
	tmpl, err := tmpl.Parse(structTmpl)
	if err != nil {
		t.Fatalf("Failed to parse struct: %v", err)
	}
	// fill out the template
	args := struct {
		Name           string
		URL            string
		IgnoreElements []string
		Fields         []Field
	}{
		Name:           "TestStruct",
		URL:            "https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
		IgnoreElements: []string{"script", "style", "link", "img", "footer", "header"},
		Fields: []Field{
			{
				Name:            "A",
				Type:            "string",
				Description:     "A description of the field",
				HeaderSelector:  "tr:nth-child(1) td:nth-child(1)",
				DataSelector:    "tr td:nth-child(1)",
				ControlSelector: "$text",
				MustBePresent:   "NCAA Codes",
			},
			{
				Name:            "B",
				Type:            "int",
				Description:     "A description of the field",
				HeaderSelector:  "tr:nth-child(1) td:nth-child(2)",
				DataSelector:    "tr td:nth-child(2)",
				ControlSelector: "$text",
				MustBePresent:   "NCAA Codes",
			},
		},
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "struct", args)
	if err != nil {
		t.Fatalf("Failed to execute struct template: %v", err)
	}
	// write the buffer to the file
	a.NotEqual(buf.String(), "")
	t.Logf("struct: %s", buf.String())
	t.Fail()
}

func TestStructTemplatePrompt(t *testing.T) {
	a := assert.New(t)
	tmpl := template.New("struct_file_template")
	tmpl, err := tmpl.Parse(structTmpl)
	if err != nil {
		t.Fatalf("Failed to parse struct: %v", err)
	}
	// fill out the template
	args := struct {
		Name           string
		URL            string
		IgnoreElements []string
		Fields         []Field
		Content        string
		Selectors      []Selector
	}{
		Name:           "TestStruct",
		URL:            "https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
		IgnoreElements: []string{"script", "style", "link", "img", "footer", "header"},
		Fields: []Field{
			{
				Name:            "A",
				Type:            "string",
				Description:     "A description of the field",
				HeaderSelector:  "tr:nth-child(1) td:nth-child(1)",
				DataSelector:    "tr td:nth-child(1)",
				ControlSelector: "$text",
				MustBePresent:   "NCAA Codes",
			},
			{
				Name:            "B",
				Type:            "int",
				Description:     "A description of the field",
				HeaderSelector:  "tr:nth-child(1) td:nth-child(2)",
				DataSelector:    "tr td:nth-child(2)",
				ControlSelector: "$text",
				MustBePresent:   "NCAA Codes",
			},
		},
	}
	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "prompt", args)
	a.NoError(err)
	// write the buffer to the file
	a.NotEqual(buf.String(), "")
	t.Logf("struct: %s", buf.String())
	t.Fail()
}
