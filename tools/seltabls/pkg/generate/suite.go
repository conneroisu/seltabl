// Package generate is a package for generating a suite of files for a given url.
package generate

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/generate/config"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/generate/struc"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/generate/test"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/sashabaranov/go-openai"
	"github.com/yosssi/gohtml"
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
		structFile := domain.StructFile{
			Name:           name,
			URL:            url,
			IgnoreElements: ignoreElements,
			ConfigFile:     configFile,
			TreeWidth:      treeWidth,
			TreeDepth:      treeDepth,
			HTMLContent:    htmlBody,
		}
		err = struc.Generate(ctx, client, &structFile, &configFile)
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
		return nil
	}
}

// writeFile writes a file to the given path
func writeFile(name string, content string) error {
	log.Debugf("Write file called with name: %s", name)
	defer log.Debugf("Write file finished with name: %s", name)
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// GetURL gets the url and returns the body of the http response.
//
// If an error occurs, it returns an error.
func GetURL(url string, ignoreElements []string) ([]byte, error) {
	log.Debugf("Get URL called with url: %s", url)
	defer log.Debugf("Get URL finished with url: %s", url)
	cli := http.DefaultClient
	resp, err := cli.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get url: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get url: %s", resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	doc, err := parsers.GetMinifiedDoc(
		string(body),
		ignoreElements,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get minified doc: %w", err)
	}
	docHTML, err := doc.Html()
	if err != nil {
		return nil, fmt.Errorf("failed to get html: %w", err)
	}
	docHTML = gohtml.FormatWithLineNo(docHTML)
	docHTML = strings.ReplaceAll(docHTML, "\n", "")
	return []byte(docHTML), nil
}
