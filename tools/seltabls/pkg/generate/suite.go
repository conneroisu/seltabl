// Package generate is a package for generating a suite of files for a given url.
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
	treeWidth int,
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
		err = configFile.Generate(ctx, client)
		if err != nil {
			return fmt.Errorf("failed to generate config file: %w", err)
		}
		structFile := StructFile{
			Name:           name,
			URL:            url,
			IgnoreElements: ignoreElements,
			ConfigFile:     configFile,
			TreeWidth:      treeWidth,
		}
		err = structFile.Generate(ctx, client)
		if err != nil {
			return fmt.Errorf("failed to generate struct file: %w", err)
		}
		testFile := TestFile{
			Name:       name,
			URL:        url,
			ConfigFile: configFile,
			StructFile: structFile,
		}
		err = testFile.Generate(ctx, client)
		if err != nil {
			return fmt.Errorf("failed to generate test file: %w", err)
		}
		return nil
	}
}
