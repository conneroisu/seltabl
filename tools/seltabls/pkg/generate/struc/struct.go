// Package struc contains the struct file generation logic.
package struc

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"os"

	// Embedded for the struct template
	_ "embed"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
	"github.com/sashabaranov/go-openai"
)

// Generate generates a struct file for the given name.
//
// If the context is cancelled, it returns an error from the context.
func Generate(
	ctx context.Context,
	client *openai.Client,
	sF *domain.StructFile,
	cfg *domain.ConfigFile,
	section *domain.Section,
) error {
	log.Debugf("Generate called on stuct: %v", sF)
	defer log.Debugf("Cenerate called on stuct: %v", sF)
	if !domain.IsValidTreeDepth(sF.TreeDepth) ||
		!domain.IsValidTreeWidth(sF.TreeWidth) {
		if domain.IsValidTreeDepth(sF.TreeDepth) {
			return fmt.Errorf("tree depth is not valid: %d", sF.TreeDepth)
		}
		if domain.IsValidTreeWidth(sF.TreeWidth) {
			return fmt.Errorf("tree width is not valid: %d", sF.TreeWidth)
		}
	}
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		f, err := os.Create(sF.Name + "_seltabl.go")
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()
		structFile, err := generate(
			ctx,
			f,
			sF,
			cfg,
			client,
			*section,
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
func generate(
	ctx context.Context,
	writer io.Writer,
	sF *domain.StructFile,
	cfg *domain.ConfigFile,
	client *openai.Client,
	section domain.Section,
) (domain.StructFile, error) {
	log.Debugf("generate called on stuct: %v", sF)
	defer log.Debugf("generate called on stuct: %v", sF)
	structStruct, err := NewStructStruct(
		sF.Name,
		sF.URL,
		sF.IgnoreElements,
		&section,
	)
	if err != nil {
		return *sF, fmt.Errorf("failed to create struct prompt: %w", err)
	}
	// write to the struct file
	_, err = writer.Write([]byte(structStruct))
	if err != nil {
		return *sF, fmt.Errorf("failed to write struct file: %w", err)
	}
	return *sF, nil
}
