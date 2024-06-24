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
		return nil, fmt.Errorf("failed to parse structs: %w", err)
	}
	for _, st := range sts {
		diags, err := s.getDiagnosticsForStruct(st, data)
		if err != nil {
			return nil, fmt.Errorf("failed to get diagnostics for struct: %w", err)
		}
		diagnostics = append(
			diagnostics,
			diags...,
		)
	}
	return diagnostics, nil
}

// getDiagnosticsForStruct returns diagnostics for a given struct
func (s *State) getDiagnosticsForStruct(
	strt parsers.Structure,
	data parsers.StructCommentData,
) (diagnostics []lsp.Diagnostic, err error) {
	content, err := s.clientGet(data.URLs[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get the content of the url: %w", err)
	}
	wg := conc.WaitGroup{}
	for j := range strt.Fields {
		for i := range strt.Fields[j].Tags.Len() {
			wg.Go(func() {
				for k := range diagnosticKeys {
					if diagnosticKeys[k] == strt.Fields[j].Tags.Tag(i).Key {
						selector := strt.Fields[j].Tags.Tag(i).Value()
						verified, err := s.validateSelector(selector, content)
						if !verified || err != nil {
							diag := lsp.Diagnostic{
								Range: lsp.LineRange(
									strt.Fields[j].Line-1,
									strt.Fields[j].Tags.Tag(i).Start,
									strt.Fields[j].Tags.Tag(i).End,
								),
								Severity: lsp.DiagnosticWarning,
								Source:   "seltabls",
							}
							if err != nil {
								diag.Message = fmt.Sprintf(
									"Failed to validate selector %s against known url content: %s",
									strt.Fields[j].Tags.Tag(i).Value(),
									err.Error(),
								)
								diagnostics = append(diagnostics, diag)
								return
							}
							diag.Message = fmt.Sprintf(
								"Could not verify selector %s against known url content",
								strt.Fields[j].Tags.Tag(i).Value(),
							)
							diagnostics = append(diagnostics, diag)
						}
					}
				}
			})
		}
	}
	wg.Wait()
	return diagnostics, nil
}

// validateSelector validates a selector against a known url content in the form of a goquery document
func (s *State) validateSelector(selector string, doc *goquery.Document) (bool, error) {
	// Create a new goquery document from the response body
	selection := doc.Find(selector)
	// Check if the selector is in the response body
	if selection.Length() < 1 {
		return false, nil
	}
	return true, nil
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
		return nil, fmt.Errorf("failed to create a new goquery document: %v", err)
	}
	return doc, nil
}
