package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// get gets the html from a given url
//
// It is used to get the html from a given url for presenting possible fields in the form.
func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to get url: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read body: %w", err)
	}
	body, err = simplifyBody(body)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to simplify body: %w", err)
	}
	return body, nil

}

// simplifyBody simplifies the body by removing all script tags, style tags,
func simplifyBody(body []byte) ([]byte, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}
	// remove all script tags
	doc.Find("script").Remove()
	// remove all style tags
	doc.Find("style").Remove()
	// remove all link tags
	doc.Find("link").Remove()
	// remove all comments
	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		s.Remove()
	})
	str, err := doc.Html()
	if err != nil {
		return nil, fmt.Errorf("failed to get html: %w", err)
	}
	return []byte(str), nil
}
