package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// DefaultClientGet gets the html of a given url.
//
// It uses the http.DefaultClient to make a GET request to the given url and
// returns the goquery document.
func DefaultClientGet(url string) (*goquery.Document, error) {
	// Http request to the server
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create a new request: %v",
			err,
		)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	done, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send the request: %v", err)
	}
	defer done.Body.Close()
	body, err := io.ReadAll(done.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read the response body: %v",
			err,
		)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create a new goquery document: %v",
			err,
		)
	}
	return doc, nil
}
