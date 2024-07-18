package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewFileTest tests the TestFilePromptArgs struct with the NewFile
// function.
func TestNewFileTest(t *testing.T) {
	a := assert.New(t)
	content, err := NewFile(
		FileTest{
			Version:     "v0.0.0",
			Name:        "TestStruct",
			URL:         "https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
			PackageName: "main",
		},
	)
	a.NoError(err)
	a.NotEmpty(content)
	t.Logf("struct: %s", content)
}

// TestNewFileStruct tests the StructFilePromptArgs struct with the NewFile
// function.
func TestNewFileStruct(t *testing.T) {
	a := assert.New(t)
	content, err := NewFile(
		FileStruct{
			PackageName: "main",
			Name:        "TestStruct",
			URL:         "https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html",
			IgnoreElements: []string{
				"script",
				"style",
				"link",
				"img",
				"footer",
				"header",
			},
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
		},
	)
	a.NoError(err)
	a.NotEmpty(content)
	t.Logf("struct: %s", content)
}
