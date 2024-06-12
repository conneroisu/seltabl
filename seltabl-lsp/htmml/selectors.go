package htmml

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// getAllSelectors retrieves all selectors from the given HTML document
func getAllSelectors(doc *goquery.Document) []string {
	selectorsMap := make(map[string]struct{})
	// Function to process each element and store its selector
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		var selectorParts []string
		if tag := goquery.NodeName(s); tag != "" {
			selectorParts = append(selectorParts, tag)
		}
		if id, exists := s.Attr("id"); exists {
			selectorParts = append(selectorParts, "#"+id)
		}
		if class, exists := s.Attr("class"); exists {
			classParts := strings.Fields(class)
			for _, classPart := range classParts {
				selectorParts = append(selectorParts, "."+classPart)
			}
		}
		selector := strings.Join(selectorParts, "")
		selectorsMap[selector] = struct{}{}
	})
	// Convert map to slice
	selectors := make([]string, 0, len(selectorsMap))
	for selector := range selectorsMap {
		selectors = append(selectors, selector)
	}
	return selectors
}
