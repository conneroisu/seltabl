package seltabl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
//	     <tr> <td>a</td> <td>b</td> </tr>
//	     <tr> <td>1</td> <td>2</td> </tr>
//	     <tr> <td>3</td> <td>4</td> </tr>
//	     <tr> <td>5</td> <td>6</td> </tr>
//	     <tr> <td>7</td> <td>8</td> </tr>
//	</table>
//	`
//
//	type FixtureStruct struct {
//	        A string `json:"a" hSel:"tr:nth-child(1)" dSel:"table tr:not(:first-child) td:nth-child(1)" cSel:"$text"`
//	        B string `json:"b" hSel:"tr:nth-child(1)" dSel:"table tr:not(:first-child) td:nth-child(2)" cSel:"$text"`
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
	if dType.Kind() != reflect.Struct && dType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("expected struct, got %s", dType.Kind())
	}
	results := make([]T, 0)
	var cfg *SelectorConfig
	for i := 0; i < dType.NumField(); i++ {
		cfg = NewSelectorConfig(dType.Field(i).Tag)
		if cfg.DataSelector == "" {
			continue
		}
		dataRows := doc.Find(cfg.DataSelector)
		if dataRows.Length() <= 0 {
			return nil, ErrSelectorNotFound{
				Typ:   dType,
				Field: dType.Field(i),
				Cfg:   cfg,
			}
		}
		if cfg.HeadSelector != "" && cfg.HeadSelector != "-" {
			_ = dataRows.RemoveFiltered(cfg.HeadSelector)
		}
		if len(results) < dataRows.Length() {
			results = make([]T, dataRows.Length())
		}
		for j := 0; j < dataRows.Length(); j++ {
			err := SetStructField(
				&results[j],
				dType.Field(i), // name of the field to set
				dataRows.Eq(j), // goquery selection for cell
				&selector{
					control: cfg.ControlTag,
					query:   cfg.QuerySelector,
				}, // selector for the inner cell
			)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to set field %s: %s",
					dType.Field(i).Name,
					err,
				)
			}
		}
	}
	if len(results) < 1 {
		return nil, fmt.Errorf("no data found")
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
//	type FixtureStruct struct {
//	        A string `json:"a" hSel:"tr:nth-child(1)" dSel:"table tr:not(:first-child) td:nth-child(1)" cSel:"$text"`
//	        B string `json:"b" hSel:"tr:nth-child(1)" dSel:"table tr:not(:first-child) td:nth-child(2)" cSel:"$text"`
//	}
//
//	func main() {
//		p, err := seltabl.NewFromString[TableStruct](`
//		<table>
//			<tr> <td>a</td> <td>b</td> </tr>
//			<tr> <td>1</td> <td>2</td> </tr>
//			<tr> <td>3</td> <td>4</td> </tr>
//			<tr> <td>5</td> <td>6</td> </tr>
//			<tr> <td>7</td> <td>8</td> </tr>
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
	reader := strings.NewReader(htmlInput)
	doc, err := goquery.NewDocumentFromReader(reader)
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
//		A string `json:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
//		B string `json:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
//	}
//
//	func main() {
//		p, err := seltabl.NewFromReader[TableStruct](strings.NewReader(`
//		<table>
//			<tr> <td>a</td> <td>b</td> </tr>
//			<tr> <td>1</td> <td>2</td> </tr>
//			<tr> <td>3</td> <td>4</td> </tr>
//			<tr> <td>5</td> <td>6</td> </tr>
//			<tr> <td>7</td> <td>8</td> </tr>
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
//	type FixtureStruct struct {
//	        A string `json:"a" hSel:"tr:nth-child(1)" dSel:"table tr:not(:first-child) td:nth-child(1)" cSel:"$text"`
//	        B string `json:"b" hSel:"tr:nth-child(1)" dSel:"table tr:not(:first-child) td:nth-child(2)" cSel:"$text"`
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
		strings.NewReader(string(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}
	return New[T](doc)
}

// NewFromBytes parses a byte slice into a slice of structs adhering to the
// given generic type.
//
// The byte slice must be a valid html page with a single table.
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
//	type FixtureStruct struct {
//	        A string `json:"a" hSel:"tr:nth-child(1)" dSel:"table tr:not(:first-child) td:nth-child(1)" cSel:"$text"`
//	        B string `json:"b" hSel:"tr:nth-child(1)" dSel:"table tr:not(:first-child) td:nth-child(2)" cSel:"$text"`
//	}
//
//	func main() {
//		p, err := seltabl.NewFromBytes[TableStruct]([]byte(`
//		<table>
//			<tr> <td>a</td> <td>b</td> </tr>
//			<tr> <td>1</td> <td>2</td> </tr>
//			<tr> <td>3</td> <td>4</td> </tr>
//			<tr> <td>5</td> <td>6</td> </tr>
//			<tr> <td>7</td> <td>8</td> </tr>
//		</table>
//		`))
//		if err != nil {
//			panic(err)
//		}
//		for _, pp := range p {
//			fmt.Printf("pp %+v\n", pp)
//		}
//	}
func NewFromBytes[T any](b []byte) ([]T, error) {
	doc, err := goquery.NewDocumentFromReader(
		strings.NewReader(string(b)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}
	return New[T](doc)
}

// NewCh parses a goquery doc into a slice of structs delivered to a channel.
//
// It parse the html for each slice of the structs.
//
// The struct given as an argument must have a field with the
// tag seltabl, a header selector with the tag hSel, and a data
// selector with the tag key dSel.
func NewCh[T any](doc *goquery.Document, ch chan T) error {
	dType := reflect.TypeOf((*T)(nil)).Elem()
	if dType.Kind() != reflect.Struct && dType.Kind() != reflect.Ptr {
		return fmt.Errorf("expected struct, got %s", dType.Kind())
	}
	results := make([]T, 0)
	var cfg *SelectorConfig
	cfg = NewSelectorConfig(dType.Field(0).Tag)
	dRows := doc.Find(cfg.DataSelector)
	for i := 0; i < dRows.Length(); i++ {
		if len(results) < dRows.Length() {
			results = make([]T, dRows.Length())
		}
		// now we get the data row of i instead of for each field
		// so apply downwards kinda only on this row
		for j := 0; j < dType.NumField(); j++ {
			cfg = NewSelectorConfig(dType.Field(j).Tag)
			if cfg.DataSelector == "" {
				continue
			}
			dataRows := doc.Find(cfg.DataSelector).Eq(i)
			if cfg.HeadSelector != "" && cfg.HeadSelector != "-" {
				_ = dataRows.RemoveFiltered(cfg.HeadSelector)
			}
			err := SetStructField(
				&results[i],
				dType.Field(j), // name of the field to set
				dataRows,       // goquery selection for cell
				&selector{
					control: cfg.ControlTag,
					query:   cfg.QuerySelector,
				}, // selector for the inner cell
			)
			if err != nil {
				return fmt.Errorf(
					"failed to set field %s: %s",
					dType.Field(j).Name,
					err,
				)
			}
		}
		ch <- results[i]
	}
	return nil
}

// NewFromReaderCh parses a reader into a slice of structs.
func NewFromReaderCh[T any](r io.Reader, ch chan T) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return fmt.Errorf("failed to parse html: %w", err)
	}
	return NewCh(doc, ch)
}

// NewFromStringCh parses a string into a slice of structs.
func NewFromStringCh[T any](htmlInput string, ch chan T) error {
	reader := strings.NewReader(
		htmlInput,
	)
	doc, err := goquery.NewDocumentFromReader(
		reader,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to create new goquery document from reader: %w",
			err,
		)
	}
	return NewCh(doc, ch)
}

// NewFromBytesCh parses a byte slice into a slice of structs adhering to the
// given generic type.
func NewFromBytesCh[T any](b []byte, ch chan T) error {
	doc, err := goquery.NewDocumentFromReader(
		strings.NewReader(
			string(b),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to parse html: %w", err)
	}
	return NewCh(doc, ch)
}

// NewFromURLCh parses a given URL's html into a slice of structs adhering to the
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
//	type FixtureStruct struct {
//	        A string `json:"a" hSel:"tr:nth-child(1)" dSel:"table tr:not(:first-child) td:nth-child(1)" cSel:"$text"`
//	        B string `json:"b" hSel:"tr:nth-child(1)" dSel:"table tr:not(:first-child) td:nth-child(2)" cSel:"$text"`
//	}
//
//	func main() {
//		p, err := seltabl.NewFromURLCh[TableStruct]("https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html", ch)
//		if err != nil {
//			panic(err)
//		}
//		for _, pp := range p {
//			fmt.Printf("pp %+v\n", pp)
//		}
//	}
func NewFromURLCh[T any](url string, ch chan T) error {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get url: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to parse html: %w", err)
	}
	return NewCh(doc, ch)
}

// NewChFn parses a reader into a channel of structs.
//
// It also applies a function to each struct before adding it to the channel.
func NewChFn[
	T any,
	F func(T) bool,
](
	doc *goquery.Document,
	ch chan T,
	fn F,
) error {
	dType := reflect.TypeOf((*T)(nil)).Elem()
	if dType.Kind() != reflect.Struct && dType.Kind() != reflect.Ptr {
		return fmt.Errorf("expected struct, got %s", dType.Kind())
	}
	results := make([]T, 0)
	var cfg *SelectorConfig
	cfg = NewSelectorConfig(dType.Field(0).Tag)
	dRows := doc.Find(cfg.DataSelector)
	for i := 0; i < dRows.Length(); i++ {
		if len(results) < dRows.Length() {
			results = make([]T, dRows.Length())
		}
		for j := 0; j < dType.NumField(); j++ {
			cfg = NewSelectorConfig(dType.Field(j).Tag)
			if cfg.DataSelector == "" {
				continue
			}
			dataRows := doc.Find(cfg.DataSelector).Eq(i)
			if cfg.HeadSelector != "" && cfg.HeadSelector != "-" {
				_ = dataRows.RemoveFiltered(cfg.HeadSelector)
			}
			err := SetStructField(
				&results[i],
				dType.Field(j), // name of the field to set
				dataRows,       // goquery selection for cell
				&selector{
					control: cfg.ControlTag,
					query:   cfg.QuerySelector,
				}, // selector for the inner cell
			)
			if err != nil {
				return fmt.Errorf(
					"failed to set field %s: %s",
					dType.Field(j).Name,
					err,
				)
			}
		}
		if fn(results[i]) {
			ch <- results[i]
		}
	}
	return nil
}

// NewFromReaderChFn parses a reader into a channel of structs.
// It also applies a function to each struct before adding it to the channel.
func NewFromReaderChFn[
	T any,
	F func(T) bool,
](
	r io.Reader,
	ch chan T,
	fn F,
) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return fmt.Errorf("failed to parse html: %w", err)
	}
	return NewChFn(doc, ch, fn)
}

// NewFromStringChFn parses a string into a channel of structs.
// It also applies a function to each struct before adding it to the channel.
func NewFromStringChFn[
	T any,
	F func(T) bool,
](
	htmlInput string,
	ch chan T,
	fn F,
) error {
	reader := strings.NewReader(htmlInput)
	doc, err := goquery.NewDocumentFromReader(
		reader,
	)
	if err != nil {
		return fmt.Errorf("failed to parse html: %w", err)
	}
	return NewChFn(doc, ch, fn)
}

// NewFromBytesChFn parses a byte slice into a channel of structs.
// It also applies a function to each struct before adding it to the channel.
func NewFromBytesChFn[
	T any,
	F func(T) bool,
](
	b []byte,
	ch chan T,
	fn F,
) error {
	doc, err := goquery.NewDocumentFromReader(
		strings.NewReader(string(b)),
	)
	if err != nil {
		return fmt.Errorf("failed to parse html: %w", err)
	}
	return NewChFn(doc, ch, fn)
}

// NewChFnErr parses a reader into a channel of structs.
//
// It also applies a function to each struct before adding it to the channel.
//
// It ignores errors per row selected.
func NewChFnErr[
	T any,
	F func(T) bool,
](
	doc *goquery.Document,
	ch chan T,
	fn F,
) error {
	dType := reflect.TypeOf((*T)(nil)).Elem()
	if dType.Kind() != reflect.Struct && dType.Kind() != reflect.Ptr {
		return fmt.Errorf("expected struct, got %s", dType.Kind())
	}
	results := make([]T, 0)
	var cfg *SelectorConfig
	cfg = NewSelectorConfig(dType.Field(0).Tag)
	dRows := doc.Find(cfg.DataSelector)
	for i := 0; i < dRows.Length(); i++ {
		if len(results) < dRows.Length() {
			results = make([]T, dRows.Length())
		}
		for j := 0; j < dType.NumField(); j++ {
			cfg = NewSelectorConfig(dType.Field(j).Tag)
			if cfg.DataSelector == "" {
				continue
			}
			dataRows := doc.Find(cfg.DataSelector).Eq(i)
			if cfg.HeadSelector != "" && cfg.HeadSelector != "-" {
				_ = dataRows.RemoveFiltered(cfg.HeadSelector)
			}
			err := SetStructField(
				&results[i],
				dType.Field(j), // name of the field to set
				dataRows,       // goquery selection for cell
				&selector{
					control: cfg.ControlTag,
					query:   cfg.QuerySelector,
				}, // selector for the inner cell
			)
			if err != nil {
				break
			}
		}
		if fn(results[i]) {
			ch <- results[i]
		}
	}
	return nil
}

// NewFromStringChFnErr parses a string into a channel of structs.
//
// It also applies a function to each struct before adding it to the channel.
//
// It ignores errors per row selected.
func NewFromStringChFnErr[T any](
	htmlInput string,
	ch chan T,
	fn func(T) bool,
) error {
	doc, err := goquery.NewDocumentFromReader(
		strings.NewReader(htmlInput),
	)
	if err != nil {
		return fmt.Errorf("failed to parse html: %w", err)
	}
	return NewChFnErr(doc, ch, fn)
}
