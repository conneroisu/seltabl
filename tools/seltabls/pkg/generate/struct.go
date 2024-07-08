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
	Name           string   `json:"-"`
	URL            string   `json:"-"`
	IgnoreElements []string `json:"ignore-elements"`

	Fields []Field `json:"fields"`
}

// Field is a struct for a field
type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
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
