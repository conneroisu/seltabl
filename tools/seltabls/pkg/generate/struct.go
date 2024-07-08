package generate

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"

	"github.com/sashabaranov/go-openai"
)

// StructFile is a struct for a struct file.
//
// It contains attributes relating to the name, url, and ignore elements of the struct file.
type StructFile struct {
	// Name is the name of the struct file.
	Name string `json:"-"`
	// URL is the url for the struct file.
	URL string `json:"-"`
	// IgnoreElements is a list of elements to ignore when generating the struct.
	IgnoreElements []string `json:"ignore-elements"`

	// Fields is a list of fields for the struct.
	Fields []Field `json:"fields"`
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
		f, err := os.Create(s.Name + ".go")
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()
		// Create a new buffer
		w := new(bytes.Buffer)
		// Create a new template
		tmpl := template.New("struct_file_template")
		// Execute the template
		err = tmpl.Execute(w, s)
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
