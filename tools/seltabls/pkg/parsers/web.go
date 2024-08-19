package parsers

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/yosssi/gohtml"
)

const (
	childsep = ">"
	empty    = ""
)

// GetAllSelectors retrieves all selectors from the given HTML document.
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
		html, _ := doc.Html()
		return nil, fmt.Errorf("no selectors found in document: %s", html)
	}
	return strs, nil
}

// getSelectorsFromSelection returns the CSS selector for the given goquery selection
func getSelectorsFromSelection(s *goquery.Selection) string {
	if s.Length() == 0 {
		return empty
	}
	// Recursive call for the parent
	parentSelector := getSelectorsFromSelection(s.Parent())
	// Get the selector for the current node
	currentSelector := singleSelector(s)
	// Combine the parent and current selectors
	if parentSelector != empty && currentSelector != "" &&
		parentSelector != currentSelector {
		return parentSelector + childsep + currentSelector
	} else if parentSelector != empty && currentSelector == "" {
		return parentSelector
	}

	return currentSelector
}

// GetSelectors gets all the selectors from the given URL and appends them to
// the selectors slice.
func GetSelectors(
	ctx context.Context,
	db *data.Database[master.Queries],
	url string,
	ignores []string,
	mustOccur int,
) (selectors []master.Selector, err error) {
	var doc *goquery.Document
	tx, err := db.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	rows, err := db.Queries.GetSelectorsByMinOccurances(
		ctx,
		master.GetSelectorsByMinOccurancesParams{
			Value:      url,
			Occurances: int64(mustOccur),
		},
	)
	if err == nil && len(rows) > 0 {
		var selectors []master.Selector
		for _, row := range rows {
			selectors = append(selectors, *row)
		}
		return selectors, nil
	}
	doc, err = GetMinifiedDoc(url, ignores)
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
	if err != nil {
		return nil, fmt.Errorf("failed to insert html: %w", err)
	}
	URL, err := db.Queries.UpsertURL(
		ctx,
		master.UpsertURLParams{Value: url, HtmlID: HTML.ID},
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
		if found.Length() == 0 {
			continue
		}
		selectorContext, err := found.Parent().First().Html()
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

		if found.Length() >= mustOccur {
			selectors = append(selectors, *selector)
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return selectors, nil
}

// singleSelector returns a single CSS selector for the given node.
func singleSelector(selection *goquery.Selection) string {
	var selector string
	id, exists := selection.Attr("id")
	if exists {
		checkingVal := strings.Join(strings.Fields(id), ".")
		if !strings.Contains(checkingVal, "content") {
			selector = fmt.Sprintf("%s#%s", goquery.NodeName(selection), id)
		}
	}
	attr, exists := selection.Attr("class")
	if exists {
		checkingVal := strings.Join(strings.Fields(attr), ".")
		if !strings.Contains(checkingVal, "content") {
			selector = fmt.Sprintf(
				"%s.%s",
				goquery.NodeName(selection),
				strings.Join(strings.Fields(attr), "."),
			)
		}
	}
	attr, exists = selection.Attr("name")
	if exists {
		checkingVal := strings.Join(strings.Fields(attr), ".")
		if !strings.Contains(checkingVal, "content") {
			selector = fmt.Sprintf(
				"%s[name='%s']",
				goquery.NodeName(selection),
				attr,
			)
		}
	}
	attr, exists = selection.Attr("type")
	if exists {
		selector = fmt.Sprintf(
			"%s[type='%s']",
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
			"%s[value='%s']",
			goquery.NodeName(selection),
			attr,
		)
	}
	attr, exists = selection.Attr("src")
	if exists {
		selector = fmt.Sprintf(
			"%s[src='%s']",
			goquery.NodeName(selection),
			attr,
		)
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// GetMinifiedDoc gets a minified goquery doc from a given url
// and returns goquery doc and error if there is an error while
// getting the doc.
func GetMinifiedDoc(
	url string,
	disallowedTags []string,
) (doc *goquery.Document, err error) {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(body)
	doc, err = goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	for _, v := range disallowedTags {
		_ = doc.Find(v).Remove()
	}
	html, err := doc.Html()
	if err != nil {
		return nil, err
	}
	doc, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// IsValidURL checks if the given URL is valid.
func IsValidURL(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	if u.Scheme != "" && u.Host != "" {
		return fmt.Errorf("invalid URL: %s", uri)
	}
	return nil
}

func IsValidFileName(fileName string) error {
	if fileName == "" {
		return fmt.Errorf("file name is required")
	}
	// Check if the first character is a letter
	firstChar, _ := utf8.DecodeRuneInString(fileName)
	if !unicode.IsLetter(firstChar) {
		return fmt.Errorf("file name must start with a letter")
	}
	// Check the rest of the characters
	for _, ch := range fileName[1:] {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return fmt.Errorf(
				"file name must only contain letters, digits, and underscores",
			)
		}
	}
	return nil
}

func IsValidPackageName(packageName string) error {
	if packageName == "" {
		return fmt.Errorf("package name cannot be empty")
	}
	// Check if the first character is a letter
	firstChar, _ := utf8.DecodeRuneInString(packageName)
	if !unicode.IsLetter(firstChar) {
		return fmt.Errorf("package name must start with a letter")
	}
	// Check the rest of the characters
	for _, ch := range packageName[1:] {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return fmt.Errorf(
				"package name must only contain letters, digits, and underscores",
			)
		}
	}
	return nil
}
