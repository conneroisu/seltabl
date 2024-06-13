package parsers

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetAllSelectors retrieves all selectors from the given HTML document
func GetAllSelectors(doc *goquery.Document) []string {
	strs := []string{}
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
		if attr, exists := s.Attr("name"); exists {
			selectorParts = append(selectorParts, "[name="+attr+"]")
		}
		if attr, exists := s.Attr("type"); exists {
			selectorParts = append(selectorParts, "[type="+attr+"]")
		}
		if attr, exists := s.Attr("placeholder"); exists {
			selectorParts = append(selectorParts, "[placeholder="+attr+"]")
		}
		if attr, exists := s.Attr("value"); exists {
			selectorParts = append(selectorParts, "[value="+attr+"]")
		}
		if attr, exists := s.Attr("src"); exists {
			selectorParts = append(selectorParts, "[src="+attr+"]")
		}
		if attr, exists := s.Attr("href"); exists {
			selectorParts = append(selectorParts, "[href="+attr+"]")
		}
		selector := strings.Join(selectorParts, "")
		strs = append(strs, selector)
	})
	return strs
}
