package generate

import (
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

// TestFile is a struct for a test file
type TestFile struct {
	Name   string
	URL    string
	Client *openai.Client
}

// Generate generates a test file for the given name
func (t *TestFile) Generate() (err error) {
	err = t.NotNil()
	if err != nil {
		return fmt.Errorf("failed to generate test file: %w", err)
	}
	return nil
}

// NotNil checks if the client is nil
func (t *TestFile) NotNil() error {
	if t.Client == nil {
		return fmt.Errorf("client is nil")
	}
	return nil
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
	content := `package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test{{.Name}}(t *testing.T) {
	t.Parallel()
	client := NewClient()
	resp, err := client.Get("{{.URL}}")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
`
	return fmt.Sprintf(content, t)
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
