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
	treeWidth, treeDepth int,
	client *openai.Client,
	fastModel, smartModel string,
	name, url, htmlBody string,
	ignoreElements []string,
	selectors []master.Selector,
) (err error) {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		if client == nil {
			return fmt.Errorf("client is nil")
		}
		configFile := ConfigFile{
			Name:           name,
			URL:            url,
			IgnoreElements: ignoreElements,
			HTMLBody:       htmlBody,
			Selectors:      selectors,
			FastModel:      fastModel,
			SmartModel:     smartModel,
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
			TreeDepth:      treeDepth,
			HTMLContent:    htmlBody,
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
