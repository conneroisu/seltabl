package generate

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/sashabaranov/go-openai"
)

// Suite generates a suite for a given name.
//
// This suite includes a config file, a struct file, and a test file.
//
// The config file lives with the struct file for later use.
func Suite(
	ctx context.Context,
	client *openai.Client,
	name string,
	url string,
	ignoreElements []string,
	htmlBody string,
	selectors []master.Selector,
) (err error) {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		configFile := ConfigFile{
			Name:           name,
			URL:            url,
			IgnoreElements: ignoreElements,
			HTMLBody:       htmlBody,
			Selectors:      selectors,
		}
		err = configFile.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate config file: %w", err)
		}
		structFile := StructFile{
			Name:           name,
			URL:            url,
			IgnoreElements: ignoreElements,
		}
		err = structFile.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate struct file: %w", err)
		}
		testFile := TestFile{
			Name:   name,
			URL:    url,
			Client: client,
		}
		err = testFile.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate test file: %w", err)
		}
		return nil
	}
}
