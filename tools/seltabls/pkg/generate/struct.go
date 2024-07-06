package generate

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/sashabaranov/go-openai"
)

// StructFile is a struct for a struct file
type StructFile struct {
	Name   string
	URL    string
	Client *openai.Client

	fields []Field
}

// Field is a struct for a field
type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Generate generates a struct file for the given name
func (s *StructFile) Generate() error {
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
