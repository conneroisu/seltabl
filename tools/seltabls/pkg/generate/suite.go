// Package generate is a package for generating a suite of files for a given url.
package generate

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/generate/config"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/generate/struc"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/generate/test"
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
	log.Debugf(
		"Suite called with name: %s, url: %s, htmlBody: %s, ignoreElements: %v, selectors: %v",
		name,
		url,
		htmlBody,
		ignoreElements,
		selectors,
	)
	defer log.Debugf(
		"Suite completed with name: %s, url: %s, htmlBody: %s, ignoreElements: %v, selectors: %v",
		name,
		url,
		htmlBody,
		ignoreElements,
		selectors,
	)
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		if client == nil {
			return fmt.Errorf("client is nil")
		}
		configFile := domain.ConfigFile{
			Name:           name,
			URL:            url,
			IgnoreElements: ignoreElements,
			HTMLBody:       htmlBody,
			Selectors:      selectors,
			FastModel:      fastModel,
			SmartModel:     smartModel,
		}
		err = config.Generate(ctx, client, &configFile)
		if err != nil {
			return fmt.Errorf("failed to generate config file: %w", err)
		}
		for _, section := range configFile.Sections {
			structFile := domain.StructFile{
				Name:           name,
				URL:            url,
				IgnoreElements: ignoreElements,
				ConfigFile:     configFile,
				TreeWidth:      treeWidth,
				TreeDepth:      treeDepth,
				HTMLContent:    htmlBody,
				Section:        section,
			}
			err = struc.Generate(
				ctx,
				client,
				&structFile,
				&configFile,
				&section,
			)
			if err != nil {
				return fmt.Errorf("failed to generate struct file: %w", err)
			}
			testFile := domain.TestFile{
				PackageName: name,
				Name:        name,
				URL:         url,
			}
			err = test.Generate(ctx, &testFile, &configFile, client)
			if err != nil {
				return fmt.Errorf("failed to generate test file: %w", err)
			}
		}
		return nil
	}
}
