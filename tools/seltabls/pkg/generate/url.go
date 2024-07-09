package generate

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/yosssi/gohtml"
)

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
