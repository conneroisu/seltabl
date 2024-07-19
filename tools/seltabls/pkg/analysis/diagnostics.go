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
	text *string,
	data parsers.StructCommentData,
) ([]lsp.Diagnostic, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
			sts, err := parsers.ParseStructs(ctx, []byte(*text))
			if err != nil {
				return nil, fmt.Errorf(
					"failed to parse structs: %w",
					err,
				)
			}
			eg := errgroup.Group{}
			var diags []lsp.Diagnostic
			for _, st := range sts {
				eg.Go(func() error {
					content, err := http.DefaultClientGet(data.URLs[0])
					if err != nil {
						return fmt.Errorf(
							"failed to get the content of the url: %w",
							err,
						)
					}
					ds, err := st.Verify(ctx, data.URLs[0], content)
					if err != nil {
						return fmt.Errorf(
							"failed to get diagnostics for struct: %w",
							err,
						)
					}
					diags = append(diags, ds...)
					return nil
				})
			}
			if err := eg.Wait(); err != nil {
				return nil, fmt.Errorf(
					"failed to get diagnostics for struct: %w",
					err,
				)
			}
			return diags, nil
		}
	}
}
