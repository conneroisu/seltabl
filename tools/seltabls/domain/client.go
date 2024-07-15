package domain

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/yosssi/gohtml"
)

// GetRuledURL gets the url and returns the body of the http response with the html
// formatted and the ignored elements removed.
func GetRuledURL(
	url string,
	ignores []string,
) (body []byte, err error) {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get url: %w", err)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read body: %w",
			err,
		)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create document from reader: %w",
			err,
		)
	}
	for _, v := range ignores {
		_ = doc.Find(v).Remove()
	}
	html, err := doc.Html()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get html: %w",
			err,
		)
	}
	html = gohtml.FormatWithLineNo(html)
	return []byte(html), nil
}
