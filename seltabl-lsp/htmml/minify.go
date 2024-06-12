package htmml

import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
	"github.com/yosssi/gohtml"
)

// cleanHTML removes all the non-essential elements from the HTML document.
var reduceRepeaters = []string{"tr", "th", "option"}

// disallowedTags is a list of tags that are not allowed in the final
var disallowedTags = []string{"script", "style", "link", "img", "footer", "header"}

// cleanHTML removes all the non-essential elements from the HTML document.
func cleanHTML(body []byte) (string, error) {
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
	m := minify.New()
	r := bytes.NewBuffer([]byte{})
	fmtd := gohtml.Format(docHTML)
	fmt.Fprintf(r, "%s", fmtd)
	_, err = r.WriteString(fmtd)
	if err != nil {
		return "", fmt.Errorf("error writing to buffer: %w", err)
	}
	output := bytes.NewBuffer([]byte{})
	err = html.Minify(m, output, r, nil)
	if err != nil {
		return "", fmt.Errorf("error minifying html: %w", err)
	}
	return output.String(), nil
}
