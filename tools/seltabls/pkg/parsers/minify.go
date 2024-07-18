package parsers

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(body)
	doc, err = goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	var sel *goquery.Selection
	for _, v := range disallowedTags {
		sel = doc.Find(v).Remove()
	}
	html, err := sel.Html()
	if err != nil {
		return nil, err
	}
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	return doc, nil
}
