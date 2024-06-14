package minify

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/yosssi/gohtml"
)

// cleanHTML removes all the non-essential elements from the HTML document.
var reduceRepeaters = []string{"tr", "th", "option"}

// cleanHTML removes all the non-essential elements from the HTML document.
func cleanHTML(body []byte, disallowedTags []string) (string, error) {
	reader := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", fmt.Errorf(
			"clean failed error creating goquery document: %w",
			err,
		)
	}
	for _, v := range disallowedTags {
		_ = doc.Find(v).Remove()
	}
	// if a element repeats more than three time, remove all but the first three
	for _, v := range reduceRepeaters {
		_ = doc.Find(v).Each(func(i int, s *goquery.Selection) {
			if s.Parent().Children().Length() > 3 {
				s.Remove()
			}
		})
	}
	bdy := doc.Find("body")
	docHTML, err := bdy.Html()
	if err != nil {
		return "", fmt.Errorf("error getting html: %w", err)
	}
	fmtd := gohtml.Format(docHTML)
	return fmtd, nil
}

func GetMinifiedDoc(
	url string,
	disallowedTags []string,
) (*goquery.Document, error) {
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
	sanitized := colly.SanitizeFileName(url)
	// create a new file
	f, err := os.Create(sanitized)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()
	content, err := cleanHTML(body, disallowedTags)
	bytes := []byte(gohtml.Format(content))
	// write the body to the file
	_, err = f.Write(bytes)
	if err != nil {
		return nil, fmt.Errorf("error writing to file: %w", err)
	}
	fmt.Println("file written n-bytes: ", len(bytes), " bytes")
	doc, err := goquery.NewDocumentFromReader(
		strings.NewReader(
			string(body),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating goquery document: %w", err)
	}
	for _, v := range disallowedTags {
		_ = doc.Find(v).Remove()
	}
	return doc, nil
}
