package cmds

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// NewVetCmd returns the vet command
func NewVetCmd(ctx context.Context, w io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vet",
		Short: "Evaluate code for common errors or invalid selectors",
		Long: `
Similar to go vet, but for seltabl.
Evaluate code for common errors or invalid selectors.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, ctx = errgroup.WithContext(ctx)

			if len(args) == 0 {
				return fmt.Errorf("no files provided")
			}

			files, err := filepath.Glob(args[0])
			if err != nil {
				return fmt.Errorf("failed to glob files: %w", err)
			}
			for _, file := range files {
				vals, err := vetFile(ctx, file)
				if err != nil {
					return fmt.Errorf("failed to vet file: %w", err)
				}
				for _, diag := range vals {
					fmt.Printf(
						"%s\n%s\n%s\n",
						diag.Message,
						diag.Severity,
						diag.Source,
					)
				}
			}
			return nil
		},
	}
	return cmd
}

// vetFile vets a file
func vetFile(ctx context.Context, file string) ([]lsp.Diagnostic, error) {
	// if the filename is not a go file, return an error
	if filepath.Ext(file) != ".go" {
		return nil, fmt.Errorf("file is not a go file")
	}
	state, err := analysis.NewState()
	if err != nil {
		return nil, fmt.Errorf("failed to create state: %w", err)
	}
	response, err := state.OpenDocument(ctx, lsp.NotificationDidOpenTextDocument{
		Notification: lsp.Notification{
			RPC:    "2.0",
			Method: "textDocument/didOpen",
		},
		Params: lsp.DidOpenTextDocumentParams{
			TextDocument: lsp.TextDocumentItem{
				URI:        file,
				Text:       string(readFile(file)),
				LanguageID: "go",
			},
		},
	})
	return response.Params.Diagnostics, nil
}

// readFile reads a file
func readFile(file string) []byte {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return b
}
