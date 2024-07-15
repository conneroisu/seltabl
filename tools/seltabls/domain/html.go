package domain

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
)

// HtmlSel is a function for getting the html of a given css selector.
func HtmlSel(doc *goquery.Document, css string) string {
	res := doc.Find(css)
	if res.Length() == 0 {
		panic("no result found")
	}
	resp, err := res.Html()
	if err != nil {
		panic(err)
	}
	return resp
}

// HtmlContains is a function for checking if a given css selector exists in the html.
func HtmlContains(doc *goquery.Document, css string) bool {
	res := doc.Find(css)
	return res.Length() > 0
}

// HtmlContainsN is a function for checking if a given css selector exists in the html.
func HtmlContainsN(doc *goquery.Document, selectors []master.Selector) bool {
	for _, selector := range selectors {
		if HtmlContains(doc, selector.Value) {
			return true
		}
	}
	return false
}

// HtmlReduce is a function for reducing a list of selectors to a selectors contained in the html.
func HtmlReduce(doc *goquery.Document, selectors []master.Selector) (sels []master.Selector) {
	for _, selector := range selectors {
		if !HtmlContains(doc, selector.Value) {
			sels = append(sels, selector)
		}
	}
	return sels
}
