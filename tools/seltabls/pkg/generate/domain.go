package generate

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/yosssi/gohtml"
)

// writeFile writes a file to the given path
func writeFile(name string, content string) error {
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
