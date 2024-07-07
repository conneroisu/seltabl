package parsers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
)

// GetMinifiedDoc gets a minified goquery doc from a given url
// and returns goquery doc and error if there is an error while
// getting the doc.
func GetMinifiedDoc(
	url string,
	disallowedTags []string,
) (doc *goquery.Document, err error) {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, domain.ErrHTTP{
			URL:        url,
			StatusCode: resp.StatusCode,
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, domain.ErrHTTPParse{
			URL:        url,
			StatusCode: resp.StatusCode,
		}
	}
	reader := bytes.NewReader(body)
	doc, err = goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, domain.ErrDocumentFromReader{
			URL:     url,
			Content: string(body),
		}
	}
	for _, v := range disallowedTags {
		_ = doc.Find(v).Remove()
	}
	return doc, nil
}
