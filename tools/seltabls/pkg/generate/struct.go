package generate

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/sashabaranov/go-openai"
)

// StructFile is a struct for a struct file.
//
// It contains attributes relating to the name, url, and ignore elements of the struct file.
type StructFile struct {
	// Name is the name of the struct file.
	Name string `json:"-" yaml:"name"`
	// URL is the url for the struct file.
	URL string `json:"-" yaml:"url"`
	// IgnoreElements is a list of elements to ignore when generating the struct.
	IgnoreElements []string `json:"ignore-elements" yaml:"ignore-elements"`
	// Fields is a list of fields for the struct.
	Fields []Field `json:"fields" yaml:"fields"`

	// TreeWidth is the width of the tree when generating the struct.
	TreeWidth int `json:"-" yaml:"tree-width"`
	// ConfigFile is the config file for the struct file.
	ConfigFile ConfigFile `json:"-" yaml:"config-file"`
	// JSONValue is the json value for the struct yaml file.
	JSONValue string `json:"-" yaml:"json-value"`

	// Db is the database for the struct file.
	Db *data.Database[master.Queries] `json:"-" yaml:"-"`
}

// Field is a struct for a field
type Field struct {
	// Name is the name of the field.
	Name string `json:"name"`
	// Type is the type of the field.
	Type string `json:"type"`
	// Description is a description of the field.
	Description string `json:"description"`
	// HeaderSelector is the header selector for the field.
	HeaderSelector string `json:"header-selector"`
	// DataSelector is the data selector for the field.
	DataSelector string `json:"data-selector"`
	// ControlSelector is the control selector for the field.
	ControlSelector string `json:"control-selector"`
	// QuerySelector is the query selector for the field.
	QuerySelector string `json:"query-selector"`
	// MustBePresent is the must be present selector for the field.
	MustBePresent string `json:"must-be-present"`
}

// Generate generates a struct file for the given name.
//
// If the context is cancelled, it returns an error from the context.
func (s *StructFile) Generate(ctx context.Context, client *openai.Client) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		if client == nil {
			return fmt.Errorf("client is nil")
		}
		// Create a new file
		f, err := os.Create(s.Name + "_seltabl.go")
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()
		structFile, err := s.generate(
			ctx,
			client,
		)
		if err != nil {
			return fmt.Errorf("failed to generate struct: %w", err)
		}
		// Create a new buffer
		w := new(bytes.Buffer)
		// Create a new template
		tmpl := template.New("struct_file_template")
		// Execute the template
		err = tmpl.ExecuteTemplate(w, "struct", structFile)
		if err != nil {
			return fmt.Errorf("failed to execute struct file template: %w", err)
		}
		// Write the buffer to the file
		_, err = f.Write(w.Bytes())
		if err != nil {
			return fmt.Errorf("failed to write struct file: %w", err)
		}
		return nil
	}
}

// generate generates the struct file.
//
// It generates the struct file by using the given url, contents, and ignore elements.
func (s *StructFile) generate(
	ctx context.Context,
	client *openai.Client,
) (StructFile, error) {
	content, err := GetURL(s.URL, s.IgnoreElements)
	if err != nil {
		return *s, fmt.Errorf("failed to get url: %w", err)
	}
	_ = string(content)
	_, err = analysis.GetSelectors(
		ctx,
		s.Db,
		s.URL,
		s.IgnoreElements,
	)
	if err != nil {
		return *s, fmt.Errorf("failed to get selectors: %w", err)
	}
	return *s, nil
}
