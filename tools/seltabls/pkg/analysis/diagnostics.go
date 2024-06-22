package analysis

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/sourcegraph/conc"
)

var (
	diagnosticKeys = []string{
		selectorDataTag.Label,
		selectorHeaderTag.Label,
		selectorQueryTag.Label,
	}
)

// GetDiagnosticsForFile returns diagnostics for a given file
// text is all the characters in the file
func (s *State) GetDiagnosticsForFile(
	text *string,
	data parsers.StructCommentData,
) (diagnostics []lsp.Diagnostic, err error) {
	ctx := context.Background()
	sts, err := parsers.ParseStructs(ctx, []byte(*text))
	if err != nil {
		return nil, err
	}
	for _, st := range sts {
		diags := s.getDiagnosticsForStruct(st, data)
		diagnostics = append(
			diagnostics,
			diags...,
		)
	}
	return diagnostics, nil
}

// getDiagnosticsForStruct returns diagnostics for a given struct
func (s *State) getDiagnosticsForStruct(
	st parsers.Structure,
	data parsers.StructCommentData,
) []lsp.Diagnostic {
	var diagnostics []lsp.Diagnostic
	content, err := s.clientGet(data.URLs[0])
	if err != nil {
		s.Logger.Printf("failed to get the content of the url: %v\n", err)
	}
	wg := conc.WaitGroup{}
	for j := range st.Fields {
		for i := range st.Fields[j].Tags.Tags() {
			j, i := j, i
			wg.Go(func() {
				for k := range diagnosticKeys {
					if diagnosticKeys[k] == st.Fields[j].Tags.Tags()[i].Key {
						selector := st.Fields[j].Tags.Tags()[i].Value()
						verified := s.validateSelector(selector, content)
						if !verified {
							diagnostics = append(diagnostics, lsp.Diagnostic{
								Range: lsp.LineRange(
									st.Fields[j].Line-1,
									st.Fields[j].Tags.Tags()[i].Start,
									st.Fields[j].Tags.Tags()[i].End,
								),
								Severity: lsp.DiagnosticWarning,
								Source:   "seltabls",
								Message: fmt.Sprintf(
									"Could not verify selector %s against url content",
									st.Fields[j].Tags.Tags()[i].Value(),
								),
							})
						}
					}
				}
			})
		}
	}
	wg.Wait()
	return diagnostics
}

// validateSelector validates a selector
func (s *State) validateSelector(selector string, doc *goquery.Document) bool {
	// Create a new goquery document from the response body
	selection := doc.Find(selector)
	content, err := selection.Html()
	if err != nil {
		s.Logger.Printf("failed to get the html of the selector: %v\n", err)
		return false
	}
	s.Logger.Printf("Selector '%s' found selecting %s\n", selector, content)
	// Check if the selector is in the response body
	if selection.Length() < 1 {
		s.Logger.Printf(
			"Selector '%s' not found in the response body\n",
			selector,
		)
		return false
	}
	return true
}

// clientValidateSelector validates a selector using a client
func (s *State) clientGet(url string) (*goquery.Document, error) {
	// Http request to the server
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// Send the request
	client := &http.Client{}
	done, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send the request: %v", err)
	}
	defer done.Body.Close()
	// Read the response body
	body, err := io.ReadAll(done.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create a new goquery document: %v",
			err,
		)
	}
	return doc, nil

}
