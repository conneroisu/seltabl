package cmds

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/spf13/cobra"
)

// NewVetCmd returns the vet command which evaluates code for common errors or invalid selectors.
//
// Similar to go vet, but for seltabl defined structs.
func NewVetCmd(ctx context.Context, w io.Writer, r io.Reader) *cobra.Command {
	return &cobra.Command{
		Use:   "vet",
		Short: "Evaluate code for common errors or invalid selectors",
		Long: `
Similar to go vet, but for seltabl.
Evaluate code for common errors or invalid selectors.
`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SetOut(w)
			cmd.SetIn(r)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Println("No files provided")
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
}

// vetFile vets a file at the given path adhering to the given context's timeout.
func vetFile(ctx context.Context, filePath string) (response []lsp.Diagnostic, err error) {
	var state analysis.State
	if filepath.Ext(filePath) != ".go" {
		return nil, fmt.Errorf("file is not a go file")
	}
	state, err = analysis.NewState()
	if err != nil {
		return nil, fmt.Errorf("failed to create state: %w", err)
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	ctn := string(content)
	data, err := parsers.ParseStructComments(ctn)
	if err != nil {
		return response, nil
	}
	diags, err := analysis.GetDiagnosticsForFile(ctx, &state, &ctn, data)
	if err != nil {
		return nil, err
	}
	return diags, nil
}

// readFile reads a file at the given path.
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
