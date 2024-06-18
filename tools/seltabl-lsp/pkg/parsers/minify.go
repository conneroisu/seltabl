package parsers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yosssi/gohtml"
)

// cleanHTML removes all the non-essential elements from the HTML document.
func cleanHTML(body []byte, disallowedTags []string) (string, error) {
	var doc *goquery.Document
	var err error
	reader := bytes.NewReader(body)
	doc, err = goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", fmt.Errorf(
			"clean failed error creating goquery document: %w",
			err,
		)
	}
	for _, v := range disallowedTags {
		_ = doc.Find(v).Remove()
	}
	docHTML, err := doc.Find("body").Html()
	if err != nil {
		return "", fmt.Errorf("error getting html: %w", err)
	}
	return gohtml.Format(docHTML), nil
}

// GetMinifiedDoc gets the minified doc from a given url
func GetMinifiedDoc(
	url string,
	disallowedTags []string,
) (*goquery.Document, error) {
	var doc *goquery.Document
	var err error
	var content string
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
	content, err = cleanHTML(body, disallowedTags)
	doc, err = goquery.NewDocumentFromReader(
		strings.NewReader(
			string(content),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating goquery document: %w", err)
	}
	return doc, nil
}
