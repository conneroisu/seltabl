package config

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v2"
)

// ReadConfigFile reads a config file and unmarshals it into the
func ReadConfigFile(name string, cfg *domain.ConfigFile) (err error) {
	log.Debugf("ReadConfigFile called with name: %s, cfg: %v", name, cfg)
	defer log.Debugf(
		"ReadConfigFile completed with name: %s, cfg: %v, err: %v",
		name,
		cfg,
		err,
	)
	var f *os.File
	var b []byte
	f, err = os.Open(name)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	b, err = io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal file: %w", err)
	}
	return nil
}

// NewConfigFile returns a new config file with the given name, url, and ignore
// elements if not provided ignore elements this function will return a default
// ignore elements within the config file struct.
func NewConfigFile(
	name, url string,
	ignoreElements []string,
) (cfg *domain.ConfigFile) {
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
	return &domain.ConfigFile{
		Name:           name,
		URL:            url,
		IgnoreElements: ignoreElements,
	}
}

// Generate generates a config file for the given name.
func Generate(
	ctx context.Context,
	_ *openai.Client,
	c *domain.ConfigFile,
) error {
	log.Debugf(
		"Generate called with name: %s, url: %s, ignoreElements: %v",
		c.Name,
		c.URL,
		c.IgnoreElements,
	)
	defer log.Debugf(
		"Generate completed with name: %s, url: %s, ignoreElements: %v",
		c.Name,
		c.URL,
		c.IgnoreElements,
	)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		path := filepath.Join(".", "seltabl.yaml")
		log.Debugf("Writing config file to path: %s", path)
		defer log.Debugf("Config file written to path: %s", path)
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		out, err := yaml.Marshal(c)
		defer f.Close()
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		_, err = f.Write(out)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		return nil
	}
}
