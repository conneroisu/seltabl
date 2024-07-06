package analysis

import (
	"context"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/http"
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
			return nil, fmt.Errorf(
				"failed to get diagnostics for struct: %w",
				err,
			)
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
) (
	diagnostics []lsp.Diagnostic,
	err error,
) {
	content, err := http.DefaultClientGet(data.URLs[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get the content of the url: %w", err)
	}
	wg := conc.WaitGroup{}
	for j := range len(strt.Fields) {
		for i := range strt.Fields[j].Tags.Len() {
			wg.Go(func() {
				for k := range diagnosticKeys {
					if diagnosticKeys[k] == strt.Fields[j].Tags.Tag(i).Key {
						verified, err := s.validateSelector(
							strt.Fields[j].Tags.Tag(i).Value(),
							content,
						)
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
func (s *State) validateSelector(
	selector string,
	doc *goquery.Document,
) (bool, error) {
	// Create a new goquery document from the response body
	selection := doc.Find(selector)
	// Check if the selector is in the response body
	if selection.Length() < 1 {
		return false, nil
	}
	return true, nil
}
