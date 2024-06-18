package analysis

import (
	"strings"

	"github.com/conneroisu/seltabl/tools/pkg/lsp"
)

// getDiagnosticsForFile returns diagnostics for a given file
func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	for row, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "VS Code") {
			idx := strings.Index(line, "VS Code")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("VS Code")),
				Severity: 1,
				Source:   "Common Sense",
				Message:  "Please make sure we use good language in this video",
			})
		}

		if strings.Contains(line, "Neovim") {
			idx := strings.Index(line, "Neovim")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("Neovim")),
				Severity: 2,
				Source:   "Common Sense",
				Message:  "Great choice :)",
			})

		}
	}

	return diagnostics
}
