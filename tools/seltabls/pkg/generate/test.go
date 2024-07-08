package generate

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

// TestFile is a struct for a test file
type TestFile struct {
	Name string
	URL  string
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
	err = os.WriteFile(t.Name, []byte(t.Content()), 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to write test file: %w", err)
	}
	return len(p), nil
}

// Content returns the content of the test file
func (t *TestFile) Content() string {
	return fmt.Sprint(`package main`)
}

// WriteTestFile writes the test file to the file system
func WriteTestFile(name string, content string) error {
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
