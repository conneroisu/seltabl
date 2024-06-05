package seltabl

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/conneroisu/seltabl/internal/errors"
)

const (
	innerTextSelector = "$text"
	attrSelector      = "@"

	headerTag         = "seltabl"
	dataSelectorTag   = "dSel"
	headerSelectorTag = "hSel"
	cellSelectorTag   = "cSel"
)

// NewFromString parses a string into a slice of structs.
//
// The struct must have a field with the tag seltabl, a header selector with
// the tag hSel, and a data selector with the tag dSel.
func NewFromString[T any](htmlInput string) ([]T, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlInput))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}
	return New[T](doc)
}

// New parses a goquery doc into a slice of structs.
//
// The struct given as an argument must have a field with the tag seltabl, a header selector with
// the tag hSel, and a data selector with the tag dSel.
//
// The selectors responsibilties:
//
//   - header selector (hSel): used to find the header row and column for the
//     field in the given struct.
//   - data selector (dSel): used to find the data column for the field in the
//     given struct.
//   - cell selector (cSel): used to find the inner text or attribute of the
//     cell.
//
// Example:
//
//	var fixture = `
//	<table>
//
//	     <tr>
//	     	<td>a</td>
//	     	<td>b</td>
//	     </tr>
//	     <tr>
//	     	<td> 1 </td>
//	     	<td>2</td>
//	     </tr>
//	     <tr>
//	     	<td>3 </td>
//	     	<td> 4</td>
//	     </tr>
//	     <tr>
//	     	<td> 5 </td>
//	     	<td> 6</td>
//	     </tr>
//	     <tr>
//	     	<td>7 </td>
//	     	<td> 8</td>
//	     </tr>
//
//	</table>
//	`
//
//	type fixtureStruct struct {
//		A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
//		B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
//	}
//
//	func main() {
//		p, err := seltabl.New[fixtureStruct](fixture)
//		if err != nil {
//			panic(err)
//		}
//		for _, pp := range p {
//			fmt.Printf("pp %+v\n", pp)
//		}
//	}
func New[T any](doc *goquery.Document) ([]T, error) {
	dType := reflect.TypeOf((*T)(nil)).Elem()
	if dType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", dType.Kind())
	}
	results := make([]T, 0)
	for i := 0; i < dType.NumField(); i++ {
		field := dType.Field(i)
		headName := field.Tag.Get(headerTag)
		if headName == "" {
			continue
		}
		headSelector := field.Tag.Get(headerSelectorTag)
		if headSelector == "" {
			return nil, &errors.SelectorNotFound[T]{
				Struct:    *new(T),
				FieldName: headName,
				Selector:  headerSelectorTag,
			}
		}
		dataSelector := field.Tag.Get(dataSelectorTag)
		if dataSelector == "" {
			return nil, &errors.SelectorNotFound[T]{
				Struct:    *new(T),
				FieldName: headName,
				Selector:  dataSelectorTag,
			}
		}
		headRow := doc.Find(headSelector)
		if headRow.Length() == 0 {
			return nil, &errors.SelectionNotFound[T]{
				Struct:         *new(T),
				FieldName:      headName,
				SelectionQuery: headSelector,
			}
		}
		dataRows := doc.Find(
			fmt.Sprintf("%s:not(%s)", dataSelector, headSelector),
		)
		if dataRows.Length() == 0 {
			return nil, &errors.SelectionNotFound[T]{
				Struct:         *new(T),
				FieldName:      headName,
				SelectionQuery: dataSelector,
			}
		}
		cellSelector := field.Tag.Get(cellSelectorTag)
		if cellSelector == "" {
			cellSelector = innerTextSelector
		}
		if len(results) == 0 {
			results = make([]T, dataRows.Length())
		}
		for j := 0; j < dataRows.Length(); j++ {
			if err := SetStructField(
				&results[j],    // result row for this data row
				field.Name,     // name of the field to set
				dataRows.Eq(j), // goquery selection for cell
				cellSelector,   // selector for the inner cell
			); err != nil {
				return nil, fmt.Errorf(
					"failed to set field %s: %s",
					field.Name,
					err,
				)
			}
		}
	}
	if len(results) <= 1 {
		return nil, fmt.Errorf("no data found for struct %s", dType.Name())
	}
	return results, nil
}

// NewFromReader parses a reader into a slice of structs.
//
// The reader must be a valid html page with a single table.
//
// The passed in generic type must be a struct with valid selectors for the
// table and data (hSel, dSel, cSel).
func NewFromReader[T any](r io.Reader) ([]T, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}
	return New[T](doc)
}

// NewFromURL parses a URL into a slice of structs.
//
// The URL must be a valid html page with a single table.
//
// The passed in generic type must be a struct with valid selectors for the
// table and data (hSel, dSel, cSel).
func NewFromURL[T any](url string) ([]T, error) {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get url: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}
	return New[T](doc)
}
