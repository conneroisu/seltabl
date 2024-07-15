package domain

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
)

// HTMLSel is a function for getting the html of a given css selector.
func HTMLSel(doc *goquery.Document, css string) string {
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

// HTMLContains is a function for checking if a given css selector exists in the html.
func HTMLContains(doc *goquery.Document, css string) bool {
	res := doc.Find(css)
	return res.Length() > 0
}

// HTMLContainsN is a function for checking if a given css selector exists in the html.
func HTMLContainsN(doc *goquery.Document, selectors []master.Selector) bool {
	for _, selector := range selectors {
		if HTMLContains(doc, selector.Value) {
			return true
		}
	}
	return false
}

// HTMLReduce is a function for reducing a list of selectors to a selectors contained in the document.
func HTMLReduce(doc *goquery.Document, selectors []master.Selector) (sels []master.Selector) {
	for _, selector := range selectors {
		if !HTMLContains(doc, selector.Value) {
			sels = append(sels, selector)
		}
	}
	return sels
}

// HTMLReduct is a function for reducing a list of selectors to a selectors contained in the html.
func HTMLReduct(doc *goquery.Document, css string) string {
	res := doc.Find(css)
	sels, err := res.Html()
	if err != nil {
		panic(err)
	}
	return sels
}
