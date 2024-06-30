package prompts

import (
	"fmt"
	"testing"
)

func TestNewErrPrompt(t *testing.T) {
	html := `<html><body><div>Hello, World!</div></body></html>`
	content := `package main

// @url: https://example.com/one
// @ignore-elements: img, link
type MyStruct struct {
	Field1 string ` + "`json:\"field1\"`" + `
	Field2 int    ` + "`json:\"field2\"`" + `
}
`
	err := fmt.Errorf("failed to parse struct")
	url := "https://example.com/one"
	prompt, err := NewErrPrompt(html, content, url, err)
	if err != nil {
		t.Errorf("NewErrPrompt() error = %v", err)
		return
	}
	expected := `Error: failed to parse struct

HTML:
<html><body><div>Hello, World!</div></body></html>

Content:
package main

// @url: https://example.com/one
// @ignore-elements: img, link
type MyStruct struct {
	Field1 string ` + "`json:\"field1\"`" + `
	Field2 int    ` + "`json:\"field2\"`" + `
}

Error: failed to parse struct
`
	t.Log(prompt)
	t.Log(expected)
}
