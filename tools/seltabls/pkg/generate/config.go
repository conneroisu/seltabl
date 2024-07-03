package generate

import (
	"fmt"
	"os"
	"path/filepath"
)

// ConfigFile is a struct for a config file
type ConfigFile struct {
	Name           string   `yaml:"name"`
	URL            string   `yaml:"url"`
	IgnoreElements []string `yaml:"ignore-elements"`
}

// Generate generates a config file for the given name.
func (c *ConfigFile) Generate() error {
	path := filepath.Join(c.Name, c.Name+"_seltabl.yaml")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf(`
name: %s
url: %s
ignore-elements:
  - script
  - style
  - link
  - img
  - footer
  - header
`, c.Name, c.URL))
	return nil
}
