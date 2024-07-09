package generate

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

// TestFile is a struct for a test file
type TestFile struct {
	// Name is the name of the test file
	Name string `json:"name" yaml:"name"`
	// URL is the url for the test file
	URL string `json:"url"  yaml:"url"`
	// ConfigFile is the config file for the test file
	ConfigFile ConfigFile `json:"-"    yaml:"config-file"`
	// StructFile is the struct file for the test file
	StructFile StructFile `json:"-"    yaml:"struct-file"`
}

// Generate generates a test file for the given name
func (t *TestFile) Generate(ctx context.Context, client *openai.Client) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		if client == nil {
			return fmt.Errorf("client is nil")
		}
		return nil
	}
}

// Write writes the test file to the file system
func (t *TestFile) Write(p []byte) (n int, err error) {
	err = os.WriteFile(
		fmt.Sprintf("%s_test.go", t.Name),
		[]byte(t.Content()),
		0644,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to write test file: %w", err)
	}
	return len(p), nil
}

// Content returns the content of the test file
func (t *TestFile) Content() string {
	return fmt.Sprint(`package main`)
}
