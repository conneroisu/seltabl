package parsers

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/yosssi/gohtml"
)

const (
	childsep = " > "
	empty    = ""
)

// GetAllSelectors retrieves all selectors from the given HTML document
func GetAllSelectors(doc *goquery.Document) ([]string, error) {
	strs := []string{}
	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		str := getSelectorsFromSelection(s)
		if str != empty {
			if !contains(strs, str) {
				strs = append(strs, str)
			}
		}
	})
	if len(strs) == 0 {
		return nil, fmt.Errorf("no selectors found in document")
	}
	return strs, nil
}

// getSelectorsFromSelection returns the CSS selector for the given goquery selection
func getSelectorsFromSelection(s *goquery.Selection) string {
	if s.Length() == 0 {
		return ""
	}
	// Get the parent node
	parent := s.Parent()
	// Recursive call for the parent
	parentSelector := getSelectorsFromSelection(parent)
	// Get the selector for the current node
	currentSelector := singleSelector(s)
	// Combine the parent and current selectors
	if parentSelector != empty && currentSelector != "" {
		return parentSelector + childsep + currentSelector
	} else if parentSelector != empty && currentSelector == "" {
		return parentSelector
	}
	return currentSelector
}

// singleSelector returns a single CSS selector for the given node
func singleSelector(selection *goquery.Selection) string {
	var selector string
	id, exists := selection.Attr("id")
	if exists {
		selector = fmt.Sprintf("%s#%s", goquery.NodeName(selection), id)
	}
	attr, exists := selection.Attr("class")
	if exists {
		selector = fmt.Sprintf(
			"%s.%s",
			goquery.NodeName(selection),
			strings.Join(strings.Fields(attr), "."),
		)
	}
	attr, exists = selection.Attr("name")
	if exists {
		selector = fmt.Sprintf(
			"%s[name=%s]",
			goquery.NodeName(selection),
			attr,
		)
	}
	attr, exists = selection.Attr("type")
	if exists {
		selector = fmt.Sprintf(
			"%s[type=%s]",
			goquery.NodeName(selection),
			attr,
		)
	}
	attr, exists = selection.Attr("placeholder")
	if exists {
		selector = fmt.Sprintf(
			"%s[placeholder=%s]",
			goquery.NodeName(selection),
			attr,
		)
	}
	attr, exists = selection.Attr("value")
	if exists {
		selector = fmt.Sprintf(
			"%s[value=%s]",
			goquery.NodeName(selection),
			attr,
		)
	}
	attr, exists = selection.Attr("src")
	if exists {
		selector = fmt.Sprintf("%s[src=%s]", goquery.NodeName(selection), attr)
	}
	_, exists = selection.Attr("href")
	if exists {
		selector = fmt.Sprintf(
			"%s[href]",
			goquery.NodeName(selection),
		)
	}
	if selector == "" {
		selector = goquery.NodeName(selection)
	}
	return selector
}

// GetSelectors gets all the selectors from the given URL and appends them to the selectors slice
func GetSelectors(
	ctx context.Context,
	db *data.Database[master.Queries],
	url string,
	ignores []string,
) (selectors []master.Selector, err error) {
	rows, err := db.Queries.GetSelectorsByURL(
		ctx,
		master.GetSelectorsByURLParams{Value: url},
	)
	if err == nil && rows != nil {
		var selectors []master.Selector
		for _, row := range rows {
			selectors = append(selectors, *row)
		}
		return selectors, nil
	}
	doc, err := GetMinifiedDoc(url, ignores)
	if err != nil {
		return nil, fmt.Errorf("failed to get minified doc: %w", err)
	}
	docHTML, err := doc.Html()
	if err != nil {
		return nil, fmt.Errorf("failed to get html: %w", err)
	}
	HTML, err := db.Queries.InsertHTML(
		ctx,
		master.InsertHTMLParams{Value: docHTML},
	)
	URL, err := db.Queries.InsertURL(
		ctx,
		master.InsertURLParams{
			Value:  url,
			HtmlID: HTML.ID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert url: %w", err)
	}
	selectorStrings, err := GetAllSelectors(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to get selectors: %w", err)
	}
	for _, selectorString := range selectorStrings {
		found := doc.Find(selectorString)
		for i := range ignores {
			_ = found.RemoveFiltered(ignores[i])
		}
		selectorContext, err := found.First().Html()
		// selectorContext, err := found.Parent().First().Html()
		if err != nil {
			return nil, fmt.Errorf("failed to get html: %w", err)
		}
		selectorContext = gohtml.Format(selectorContext)
		selector, err := db.Queries.InsertSelector(
			ctx,
			master.InsertSelectorParams{
				Value:      selectorString,
				UrlID:      URL.ID,
				Context:    selectorContext,
				Occurances: int64(found.Length()),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert selector: %w", err)
		}
		selectors = append(selectors, *selector)
	}
	return selectors, nil
}
