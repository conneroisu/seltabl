package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConfigFile is a struct for a config file
type ConfigFile struct {
	Name           string   `yaml:"name"`
	URL            string   `yaml:"url"`
	IgnoreElements []string `yaml:"ignore-elements"`
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
func (c *ConfigFile) Generate() error {
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
	return nil
}
