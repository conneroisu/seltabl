package analysis

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"go.lsp.dev/protocol"
	"golang.org/x/sync/errgroup"
)

// GetDiagnosticsForFile returns diagnostics for a given file
// text is all the characters in the file
func GetDiagnosticsForFile(
	ctx context.Context,
	text *string,
	data *parsers.StructCommentData,
	db *data.Database[master.Queries],
) ([]protocol.Diagnostic, error) {
	log.Debug("getting diagnostics for file")
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
		var diags []protocol.Diagnostic
		log.Debugf("getting htmls for url: %s", data.URLs[0])
		html, err := db.Queries.GetHTMLByURL(ctx, master.GetHTMLByURLParams{Value: data.URLs[0]})
		if err != nil {
			return nil, fmt.Errorf("failed to get html: %w", err)
		}
		log.Debugf("creating goquery doc from html: %s", data.URLs[0])
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html.Value))
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get html: %w",
				err,
			)
		}
		for _, st := range sts {
			eg.Go(func() error {
				log.Debugf("getting diagnostics for struct: %v", st)
				ds, err := st.Verify(ctx, data.URLs[0], doc)
				if err != nil {
					return fmt.Errorf(
						"failed to get diagnostics for struct: %w",
						err,
					)
				}
				log.Debugf("got diagnostics for struct: %v", ds)
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
