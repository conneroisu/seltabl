package test

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
	"github.com/sashabaranov/go-openai"
)

// Generate generates a test file for the given name
func Generate(
	ctx context.Context,
	t *domain.TestFile,
	cfg *domain.ConfigFile,
	_ *openai.Client,
) error {
	log.Debugf("generating test file: %s", t.Name)
	defer log.Debugf("generated test file: %s", t.Name)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		version := ctx.Value("version").(string)
		ctn, err := NewTestFileContent(t.Name+"_test.go", t.URL, version, t.Name)
		if err != nil {
			return fmt.Errorf("failed to create test file content: %w", err)
		}
		_, err = t.WriteFile([]byte(ctn))
		if err != nil {
			return fmt.Errorf("failed to write test file: %w", err)
		}
		return nil
	}
}
