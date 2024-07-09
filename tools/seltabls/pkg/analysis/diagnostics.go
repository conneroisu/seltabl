package analysis

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/http"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"golang.org/x/sync/errgroup"
)

// GetDiagnosticsForFile returns diagnostics for a given file
// text is all the characters in the file
func GetDiagnosticsForFile(
	ctx context.Context,
	_ *State,
	text *string,
	data parsers.StructCommentData,
) (diagnostics []lsp.Diagnostic, err error) {
	for {
		select {
		case <-ctx.Done():
			return diagnostics, nil
		default:
			sts, err := parsers.ParseStructs(ctx, []byte(*text))
			if err != nil {
				return diagnostics, fmt.Errorf(
					"failed to parse structs: %w",
					err,
				)
			}
			eg := errgroup.Group{}
			for _, st := range sts {
				eg.Go(func() error {
					content, err := http.DefaultClientGet(data.URLs[0])
					if err != nil {
						return fmt.Errorf(
							"failed to get the content of the url: %w",
							err,
						)
					}
					diags, err := st.Verify(ctx, data.URLs[0], content)
					if err != nil {
						return fmt.Errorf(
							"failed to get diagnostics for struct: %w",
							err,
						)
					}
					diagnostics = append(
						diagnostics,
						diags...,
					)
					return nil
				})
			}
			if err := eg.Wait(); err != nil {
				return diagnostics, fmt.Errorf(
					"failed to get diagnostics for struct: %w",
					err,
				)
			}
			return diagnostics, nil
		}
	}
}
