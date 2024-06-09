package seltabl

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	cSels = []string{cSelInnerTextSelector, cSelAttrSelector}
)

const (
	cSelInnerTextSelector = "$text"   // cSelInnerTextSelector is the selector used to extract text from a cell.
	cSelAttrSelector      = "$query"  // cSelAttrSelector is the selector used to extract attributes from a cell.
	headerTag             = "seltabl" // headerTag is the tag used to mark a header cell.
	selectorDataTag       = "dSel"    // selectorDataTag is the tag used to mark a data cell.
	selectorHeaderTag     = "hSel"    // selectorHeaderTag is the tag used to mark a header selector.
	selectorControlTag    = "cSel"    // selectorControlTag is the tag used to mark a data selector.
	selectorQueryTag      = "qSel"    // selectorTag is the tag used to mark a selector.
)

// New parses a goquery doc into a slice of structs.
//
// The struct given as an argument must have a field with the
// tag seltabl, a header selector with the tag hSel, and a data
// selector with the tag dSel.
//
// The selectors responsibilities:
//
//   - header selector (hSel): used to find the header row and column for the
//     field in the given struct.
//   - data selector (dSel): used to find the data column for the field in the
//     given struct.
//   - query selector (qSel): used to query for the inner text or attribute of
//     the cell.
//   - control selector (cSel): used to control what to query for the inner
//     text or attribute of the cell.
//
// Example:
//
//	package main
//
//	var fixture = `
//	<table>
//
//	     <tr>
//	     	<td>a</td>
//	     	<td>b</td>
//	     </tr>
//	     <tr>
//	     	<td>1</td>
//	     	<td>2</td>
//	     </tr>
//	     <tr>
//	     	<td>3</td>
//	     	<td>4</td>
//	     </tr>
//	     <tr>
//	     	<td>5</td>
//	     	<td>6</td>
//	     </tr>
//	     <tr>
//	     	<td>7</td>
//	     	<td>8</td>
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
	results := make([]T, 0)
	dType := reflect.TypeOf((*T)(nil)).Elem()

	if dType.Kind() != reflect.Struct && dType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("expected struct, got %s", dType.Kind())
	}
	var cfg *SelectorConfig
	for i := 0; i < dType.NumField(); i++ {
		field := dType.Field(i)

		cfg = NewSelectorConfig(field.Tag)
		if cfg.HeadName == "" {
			continue
		}
		if cfg.HeadSelector == "" {
			return nil, fmt.Errorf(
				"header selector not found for field %s with type %s",
				field.Name,
				field.Type,
			)
		}
		if cfg.DataSelector == "" {
			return nil, fmt.Errorf(
				"data/column selector not found for field %s with type %s",
				field.Name,
				field.Type,
			)
		}
		if cfg.QuerySelector == "" || cfg.DataSelector == cSelAttrSelector {
			cfg.QuerySelector, cfg.ControlTag = cSelInnerTextSelector, cSelInnerTextSelector
		}

		dataRows := doc.Find(cfg.DataSelector)
		if dataRows.Length() == 0 {
			docHTML, err := doc.Html()
			if err != nil {
				return nil, fmt.Errorf("failed to get data rows html: %w", err)
			}
			return nil, fmt.Errorf(
				"%s found no data rows found for field %s with type %s in html: %s",
				cfg.DataSelector,
				field.Name,
				field.Type,
				docHTML,
			)
		}
		headerRow := doc.Find(cfg.HeadSelector)
		if headerRow.Length() == 0 {
			docHTML, err := doc.Html()
			if err != nil {
				return nil, fmt.Errorf("failed to get data rows html: %w", err)
			}
			return nil, fmt.Errorf(
				"%s found no header row found for field %s with type %s in html: %s",
				cfg.HeadSelector,
				field.Name,
				field.Type,
				docHTML,
			)
		}
		headRow := doc.Find(cfg.HeadSelector)
		if headRow.Length() == 0 {
			return nil, fmt.Errorf(
				"%s found no header row found for field %s with type %s",
				cfg.HeadSelector,
				field.Name,
				field.Type,
			)
		}

		if len(results) == 0 {
			results = make([]T, dataRows.Length())
		}
		for j := 0; j < dataRows.Length(); j++ {
			if err := SetStructField(
				&results[j],    // result row for this data row
				field.Name,     // name of the field to set
				dataRows.Eq(j), // goquery selection for cell
				&selector{cfg.ControlTag, cfg.QuerySelector}, // selector for the inner cell
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

// NewFromString parses a string into a slice of structs.
//
// The struct must have a field with the tag seltabl, a header selector with
// the tag hSel, and a data selector with the tag dSel.
//
// The selectors responsibilities:
//
//   - header selector (hSel): used to find the header row and column for the
//     field in the given struct.
//   - data selector (dSel): used to find the data column for the field in the
//     given struct.
//   - query selector (qSel): used to query for the inner text or attribute of
//     the cell.
//   - control selector (cSel): used to control what to query for the inner
//     text or attribute of the cell.
//
// Example:
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/conneroisu/seltabl"
//	)
//
//	type TableStruct struct {
//		A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
//		B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
//	}
//
//	func main() {
//		p, err := seltabl.NewFromString[TableStruct](`
//		<table>
//			<tr>
//				<td>a</td>
//				<td>b</td>
//			</tr>
//			<tr>
//				<td>1</td>
//				<td>2</td>
//			</tr>
//			<tr>
//				<td>3</td>
//				<td>4</td>
//			</tr>
//			<tr>
//				<td>5</td>
//				<td>6</td>
//			</tr>
//			<tr>
//				<td>7</td>
//				<td>8</td>
//			</tr>
//		</table>
//		`)
//		if err != nil {
//			panic(err)
//		}
//		for _, pp := range p {
//			fmt.Printf("pp %+v\n", pp)
//		}
//	}
func NewFromString[T any](htmlInput string) ([]T, error) {
	reader := strings.NewReader(
		htmlInput,
	)
	doc, err := goquery.NewDocumentFromReader(
		reader,
	)
	if err != nil {
		return nil,
			fmt.Errorf(
				"failed to create new goquery document from reader: %w",
				err,
			)
	}
	return New[T](doc)
}

// NewFromReader parses a reader into a slice of structs.
//
// The reader must be a valid html page with a single table.
//
// The passed in generic type must be a struct with valid selectors for the
// table and data (hSel, dSel, cSel).
//
// The selectors responsibilities:
//
//   - header selector (hSel): used to find the header row and column for the
//     field in the given struct.
//   - data selector (dSel): used to find the data column for the field in the
//     given struct.
//   - query selector (qSel): used to query for the inner text or attribute of
//     the cell.
//   - control selector (cSel): used to control what to query for the inner
//     text or attribute of the cell.
//
// Example:
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/conneroisu/seltabl"
//	)
//
//	type TableStruct struct {
//		A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
//		B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
//	}
//
//	func main() {
//		p, err := seltabl.NewFromReader[TableStruct](strings.NewReader(`
//		<table>
//			<tr>
//				<td>a</td>
//				<td>b</td>
//			</tr>
//			<tr>
//				<td>1</td>
//				<td>2</td>
//			</tr>
//			<tr>
//				<td>3</td>
//				<td>4</td>
//			</tr>
//			<tr>
//				<td>5</td>
//				<td>6</td>
//			</tr>
//			<tr>
//				<td>7</td>
//				<td>8</td>
//			</tr>
//		</table>
//		`))
//		if err != nil {
//			panic(err)
//		}
//		for _, pp := range p {
//			fmt.Printf("pp %+v\n", pp)
//		}
//	}
func NewFromReader[T any](r io.Reader) ([]T, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}
	return New[T](doc)
}

// NewFromURL parses a given URL's html into a slice of structs adhering to the
// given generic type.
//
// The URL must be a valid html page with a single table.
//
// The passed in generic type must be a struct with valid selectors for the
// table and data (hSel, dSel, cSel).
//
// The selectors responsibilities:
//
//   - header selector (hSel): used to find the header row and column for the
//     field in the given struct.
//   - data selector (dSel): used to find the data column for the field in the
//     given struct.
//   - query selector (qSel): used to query for the inner text or attribute of
//     the cell.
//   - control selector (cSel): used to control what to query for the inner
//     text or attribute of the cell.
//
// Example:
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/conneroisu/seltabl"
//	)
//
//	type TableStruct struct {
//		A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
//		B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
//	}
//
//	func main() {
//		p, err := seltabl.NewFromURL[TableStruct]("https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html")
//		if err != nil {
//			panic(err)
//		}
//		for _, pp := range p {
//			fmt.Printf("pp %+v\n", pp)
//		}
//	}
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
	doc, err := goquery.NewDocumentFromReader(
		strings.NewReader(
			string(body),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}
	return New[T](doc)
}