package generate

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v2"
)

// ConfigFile is a struct for a config file
type ConfigFile struct {
	Name           string            `yaml:"name"`
	URL            string            `yaml:"url"`
	IgnoreElements []string          `yaml:"ignore-elements"`
	HTMLBody       string            `yaml:"html-body"`
	Selectors      []master.Selector `yaml:"selectors"`
}

// ReadConfigFile reads a config file and unmarshals it into the
func (c *ConfigFile) ReadConfigFile(name string, cfg *ConfigFile) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal file: %w", err)
	}
	*cfg = *c
	return nil
}

// NewConfigFile returns a new config file with the given name, url, and ignore
// elements if not provided ignore elements this function will return a default
// ignore elements within the config file struct.
func (c *ConfigFile) NewConfigFile(
	name, url string,
	ignoreElements []string,
) *ConfigFile {
	if ignoreElements == nil {
		ignoreElements = []string{
			"script",
			"style",
			"link",
			"img",
			"footer",
			"header",
		}
	}
	return &ConfigFile{
		Name:           name,
		URL:            url,
		IgnoreElements: ignoreElements,
	}
}

// Generate generates a config file for the given name.
func (c *ConfigFile) Generate(
	ctx context.Context,
	_ *openai.Client,
) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:

		path := filepath.Join(c.Name, c.Name+"_seltabl.yaml")
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()
		_, err = f.WriteString(fmt.Sprintf(`name: %s
url: %s
ignore-elements:`+strings.Join(c.IgnoreElements, "\n  - ")+`
`, c.Name, c.URL))
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		return nil
	}
}
