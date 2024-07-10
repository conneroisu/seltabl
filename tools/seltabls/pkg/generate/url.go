package generate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/yosssi/gohtml"
)

// GetURL gets the url and returns the body of the http response.
//
// If an error occurs, it returns an error.
func GetURL(url string, ignoreElements []string) ([]byte, error) {
	log.Debugf("Get URL called with url: %s and ignoreElements: %v", url, ignoreElements)
	defer log.Debugf("Get URL finished with url: %s and ignoreElements: %v", url, ignoreElements)
	doc, err := parsers.GetMinifiedDoc(
		url,
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
