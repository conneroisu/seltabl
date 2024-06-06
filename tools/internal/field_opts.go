package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/huh"
	"github.com/conneroisu/seltabl/tools/internal/config"
)

// GetFieldOptions returns the field options for the form
//
// It is used by the GenerateCmd to generate the field options for the form
// from the config file.
func GetFieldOptions(ctf *config.Config, body []byte) ([]huh.Option[string], error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}
	tables := doc.Find("table")
	var tableHeaders []string
	for i := 0; i < tables.Length(); i++ {
		tb := tables.Eq(i)
		trs := tb.Find("tr")
		length := trs.Length()
		if length == 0 {
			fmt.Println("no trs")
			continue
		}
		for j := 0; j < 2; j++ {
			tr := trs.Eq(j)
			tr.Find("td").Each(func(_ int, s *goquery.Selection) {
				txt := strings.TrimSpace(s.Text())
				txt = strings.ReplaceAll(txt, "\n", "")
				if strings.Contains(txt, "|") {
					return
				}
				if strings.Contains(txt, "+") {
					return
				}
				_, err := strconv.ParseFloat(txt, 64)
				if err == nil {
					return
				}
				if len(txt) > 3 {
					tableHeaders = append(tableHeaders, txt)
				}
			})
			tr.Find("th").Each(func(_ int, s *goquery.Selection) {
				txt := strings.TrimSpace(s.Text())
				txt = strings.ReplaceAll(txt, "\n", "")
				if strings.Contains(txt, "|") {
					return
				}
				if strings.Contains(txt, "+") {
					return
				}
				_, err := strconv.ParseFloat(txt, 64)
				if err == nil {
					return
				}
				if len(txt) < 1 {
					return
				}
				tableHeaders = append(tableHeaders, txt)
			})
		}
	}
	opts := make([]huh.Option[string], len(tableHeaders))
	for _, tableHeader := range tableHeaders {
		opts = append(opts, huh.Option[string]{
			Key:   tableHeader,
			Value: tableHeader,
		})
	}
	return opts, nil
}
