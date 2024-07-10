// Package generate is a package for generating a suite of files for a given url.
package generate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			if client == nil {
				return fmt.Errorf("client is nil")
			}

			path := filepath.Join(".", "seltabl.yaml")
			log.Debugf("Writing config file to path: %s", path)
			defer log.Debugf("Config file written to path: %s", path)
			cfgF, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			defer cfgF.Close()
			sections, err := config.NewSections(
				ctx,
				client,
				url,
				htmlBody,
				selectors,
				fastModel,
				smartModel,
			)
			if err != nil {
				return fmt.Errorf("failed to create sections from the url's response: %w", err)
			}
			configFile, err := config.Generate(
				ctx,
				cfgF,
				&domain.ConfigFile{
					Name:           name,
					URL:            url,
					IgnoreElements: ignoreElements,
					HTMLBody:       htmlBody,
					Selectors:      selectors,
					FastModel:      fastModel,
					SmartModel:     smartModel,
					Sections:       sections,
				},
			)
			if err != nil {
				return fmt.Errorf("failed to generate config file: %w", err)
			}
			for _, section := range configFile.Sections {
				var structFile *domain.StructFile
				structFile, err = struc.Generate(
					ctx,
					client,
					&domain.StructFile{
						Name:           name,
						URL:            url,
						IgnoreElements: ignoreElements,
						ConfigFile:     configFile,
						TreeWidth:      treeWidth,
						TreeDepth:      treeDepth,
						HTMLContent:    htmlBody,
						Section:        section,
					},
					configFile,
					&section,
				)
				if err != nil {
					return fmt.Errorf("failed to generate struct file: %w", err)
				}
				var testFile *domain.TestFile
				testFile, err = test.Generate(
					ctx,
					&domain.TestFile{
						PackageName: name,
						Name:        name,
						URL:         url,
						ConfigFile:  configFile,
						StructFile:  structFile,
						Section:     &section,
					},
				)
				println(testFile.Name)
				if err != nil {
					return fmt.Errorf("failed to generate test file: %w", err)
				}
			}
			return nil
		}
	}
}
