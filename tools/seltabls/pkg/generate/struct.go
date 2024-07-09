package generate

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"os"

	_ "embed"

	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/sashabaranov/go-openai"
)

// StructFile is a struct for a struct file.
//
// It contains attributes relating to the name, url, and ignore elements of the struct file.
type StructFile struct {
	// Name is the name of the struct file.
	Name string `json:"-"               yaml:"name"`
	// URL is the url for the struct file.
	URL string `json:"-"               yaml:"url"`
	// IgnoreElements is a list of elements to ignore when generating the struct.
	IgnoreElements []string `json:"ignore-elements" yaml:"ignore-elements"`
	// Fields is a list of fields for the struct.
	Fields []Field `json:"fields"          yaml:"fields"`

	// TreeWidth is the width of the tree when generating the struct.
	TreeWidth int `json:"-" yaml:"tree-width"`
	// TreeDepth is the depth of the tree when generating the struct.
	TreeDepth int `json:"-" yaml:"tree-depth"`

	// ConfigFile is the config file for the struct file.
	ConfigFile ConfigFile `json:"-" yaml:"config-file"`
	// JSONValue is the json value for the struct yaml file.
	JSONValue string `json:"-" yaml:"json-value"`
	// HTMLContent is the html content for the struct file.
	HTMLContent string `json:"-" yaml:"html-content"`

	// Db is the database for the struct file.
	Db *data.Database[master.Queries] `json:"-" yaml:"-"`
}

//go:embed struct.tmpl
var structTmpl string

// Generate generates a struct file for the given name.
//
// If the context is cancelled, it returns an error from the context.
func (s *StructFile) Generate(
	ctx context.Context,
	client *openai.Client,
) error {
	if !isValidTreeDepth(s.TreeDepth) || !isValidTreeWidth(s.TreeWidth) {
		if isValidTreeDepth(s.TreeDepth) {
			return fmt.Errorf("tree depth is not valid: %d", s.TreeDepth)
		}
		if isValidTreeWidth(s.TreeWidth) {
			return fmt.Errorf("tree width is not valid: %d", s.TreeWidth)
		}
	}
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		f, err := os.Create(s.Name + "_seltabl.go")
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()
		identity, err := generateIdentity(ctx, s, client)
		if err != nil {
			return fmt.Errorf("failed to generate identity: %w", err)
		}
		structFile, err := s.generate(
			ctx,
			f,
			client,
			identity,
			[]Field{
				{
					Name:            "A",
					Type:            "string",
					Description:     "A description of the field",
					HeaderSelector:  "tr:nth-child(1) td:nth-child(1)",
					DataSelector:    "tr td:nth-child(1)",
					ControlSelector: "$text",
					MustBePresent:   "NCAA Codes",
				},
			},
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
			return fmt.Errorf(
				"failed to execute struct file template: %w",
				err,
			)
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
	writer io.Writer,
	client *openai.Client,
	identity IdentifyResponse,
	fields []Field,
) (StructFile, error) {
	structStruct, err := NewStructStruct(
		s.Name,
		s.URL,
		s.IgnoreElements,
		fields,
	)
	if err != nil {
		return *s, fmt.Errorf("failed to create struct prompt: %w", err)
	}
	// write to the struct file
	_, err = writer.Write([]byte(structStruct))
	if err != nil {
		return *s, fmt.Errorf("failed to write struct file: %w", err)
	}
	return *s, nil
}
