package domain

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed gen.tmpl
var genTmpl string

// GenTemplate is the must compile template
var GenTemplate = template.Must(template.New("templates").Parse(genTmpl))

type filer interface {
	template() string
}

// NewFile creates a new file for the given args.
func NewFile(
	prompt filer,
) (string, error) {
	name := prompt.template()
	buf := new(bytes.Buffer)
	err := GenTemplate.ExecuteTemplate(buf, name, prompt)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// FileStruct is the arguments for the struct file prompt.
type FileStruct struct {
	Name           string   `json:"name,omitempty"`
	URL            string   `json:"url,omitempty"`
	IgnoreElements []string `json:"ignore-elements,omitempty"`
	Fields         []Field  `json:"fields"`
	PackageName    string   `json:"package-name"`
}

func (a FileStruct) template() string { return "struct_file" }

// FileTest is the arguments for the test file prompt.
type FileTest struct {
	Version     string `json:"version"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	PackageName string `json:"package-name"`
}

func (a FileTest) template() string { return "test_file" }
