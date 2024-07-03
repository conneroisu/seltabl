package generate

import "github.com/sashabaranov/go-openai"

// StructFile is a struct for a struct file
type StructFile struct {
	Name   string
	URL    string
	Client *openai.Client
}

// GenerateStructFile generates a struct file for the given name
func (s *StructFile) Generate() error {
	return nil
}
