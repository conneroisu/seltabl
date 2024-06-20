package analysis

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/parsers"
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
func (s *State) GetDiagnosticsForFile(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	ctx := context.Background()
	urls, err := parsers.ParseStructComments(text)
	if err != nil {
		return diagnostics
	}
	sts, err := parsers.ParseStructs(ctx, []byte(text))
	if err != nil {
		return diagnostics
	}
	for _, st := range sts {
		diagnostics = append(diagnostics, s.getDiagnosticsForStruct(st, urls)...)
	}
	return diagnostics
}

// getDiagnosticsForStruct returns diagnostics for a given struct
func (s *State) getDiagnosticsForStruct(st parsers.Structure, data parsers.StructCommentData) []lsp.Diagnostic {
	var diagnostics []lsp.Diagnostic
	for _, field := range st.Fields {
		line := field.Line
		tags := field.Tags
		for _, tag := range tags.Tags() {
			for _, key := range diagnosticKeys {
				if key == tag.Key {
					verified := s.clientValidateSelector(tag.Value(), data.URLs[0])
					if !verified {
						diagnostics = append(diagnostics, lsp.Diagnostic{
							Range:    lsp.LineRange(line-1, tag.Start, tag.End),
							Severity: lsp.DiagnosticWarning,
							Source:   "seltabl-lsp",
							Message:  fmt.Sprintf("Could not verify selector %s", tag.Value()),
						})
					}
				}
			}
		}
	}
	return diagnostics
}

// validateSelector validates a selector
func (s *State) validateSelector(selector, text string) bool {
	// Create a new goquery document from the response body
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(text)))
	if err != nil {
		s.Logger.Printf("failed to create a new goquery document: %v\n", err)
		return false
	}
	// Check if the selector is in the response body
	if doc.Find(selector).Length() < 1 {
		fmt.Printf("Selector '%s' not found in the response body\n", selector)
		return false
	}
	return true
}

// clientValidateSelector validates a selector using a client
func (s *State) clientValidateSelector(selector, url string) bool {
	// Http request to the server
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.Logger.Printf("failed to create a new request: %v\n", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.Logger.Printf("failed to send the request: %v\n", err)
		return false
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.Logger.Printf("failed to read the response body: %v\n", err)
		return false
	}
	return s.validateSelector(selector, string(body))
}
