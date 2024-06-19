package parsers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

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
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	reader := bytes.NewReader(body)
	doc, err = goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("error creating goquery document: %w", err)
	}
	for _, v := range disallowedTags {
		_ = doc.Find(v).Remove()
	}
	return doc, nil
}
