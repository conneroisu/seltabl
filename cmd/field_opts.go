package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/huh"
)

// getFieldOptions returns the field options for the form
func getFieldOptions(ctf *SeltablConfig) []huh.Option[string] {
	// get request to the url
	client := &http.Client{}
	resp, err := client.Get(ctf.URL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil
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
	return opts
}
