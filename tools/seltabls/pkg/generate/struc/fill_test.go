package struc

import (
	"bytes"
	"testing"
	"text/template"

	_ "embed"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
	"github.com/stretchr/testify/assert"
)

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
		PackageName    string
		IgnoreElements []string
		Fields         []domain.Field
	}{
		Name:        "TestStruct",
		URL:         "https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
		PackageName: "main",
		IgnoreElements: []string{
			"script",
			"style",
			"link",
			"img",
			"footer",
			"header",
		},
		Fields: []domain.Field{
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
		Fields         []domain.Field
		Content        string
		Selectors      []master.Selector
	}{
		Name: "TestStruct",
		URL:  "https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
		IgnoreElements: []string{
			"script",
			"style",
			"link",
			"img",
			"footer",
			"header",
		},
		Fields: []domain.Field{
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
}

// TestStructTemplateAggregatePrompt tests the struct file template for the aggregate prompt
func TestStructTemplateAggregatePrompt(t *testing.T) {
	a := assert.New(t)
	tmpl := template.New("struct_file_template")
	tmpl, err := tmpl.Parse(structTmpl)
	// fill out the template
	args := struct {
		Name           string
		URL            string
		IgnoreElements []string
		Fields         []domain.Field
		Content        string
		Selectors      []master.Selector
		Schemas        []string
	}{
		Name:    "TestStruct",
		URL:     "https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
		Content: "<html><body><table><tr><td>a</td><td>b</td></tr><tr><td>1</td><td>2</td></tr></table></body></html>",
		IgnoreElements: []string{
			"script",
			"style",
			"link",
			"img",
			"footer",
			"header",
		},
		Schemas: []string{
			`{
				"fields": [
					{
						"name": "A",
						"type": "string",
						"description": "A description of the field",
						"header-selector": "tr:nth-child(1) td:nth-child(1)",
						"data-selector": "tr td:nth-child(1)",
						"control-selector": "$text",
						"must-be-present": "NCAA Codes"
					},
					{
						"name": "B",
						"type": "int",
						"description": "A description of the field",
						"header-selector": "tr:nth-child(1) td:nth-child(2)",
						"data-selector": "tr td:nth-child(2)",
						"control-selector": "$text",
						"must-be-present": "NCAA Codes"
					}
				]
			}`,
		},
		Fields: []domain.Field{
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
	err = tmpl.ExecuteTemplate(&buf, "aggregate", args)
	_ = a.NoError(err)
	// write the buffer to the file
	_ = a.NotEqual(buf.String(), "")
	t.Logf("struct: %s", buf.String())
}
