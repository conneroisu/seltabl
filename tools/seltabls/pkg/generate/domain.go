package generate

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/sashabaranov/go-openai"
	"github.com/yosssi/gohtml"
)

// writeFile writes a file to the given path
func writeFile(name string, content string) error {
	f, err := os.Create(name + ".go")
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

// Generatable is an interface for a generatable file.
type Generatable interface {
	Generate(ctx context.Context, client *openai.Client) error
}

// GetURL gets the url and returns the body
func GetURL(url string, ignoreElements []string) ([]byte, error) {
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
	doc, err := parsers.GetMinifiedDoc(string(body), ignoreElements)
	if err != nil {
		return nil, fmt.Errorf("failed to get minified doc: %w", err)
	}
	docHTML, err := doc.Html()
	if err != nil {
		return nil, fmt.Errorf("failed to get html: %w", err)
	}
	docHTML = gohtml.Format(docHTML)
	return []byte(docHTML), nil
}

// IsURL checks if the string is a valid URL
func IsURL(s string) error {
	_, err := url.ParseRequestURI(s)
	return err
}
