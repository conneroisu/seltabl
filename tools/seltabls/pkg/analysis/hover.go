package analysis

import (
	"bytes"
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/safe"
	"github.com/yosssi/gohtml"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
	"golang.org/x/sync/errgroup"
)

// NewHoverResponse returns a hover response for the given uri and position
func NewHoverResponse(
	ctx context.Context,
	req lsp.HoverRequest,
	db *data.Database[master.Queries],
	documents *safe.Map[uri.URI, string],
	urls *safe.Map[uri.URI, []string],
) (response *lsp.HoverResponse, err error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		response = &lsp.HoverResponse{
			Response: lsp.Response{
				RPC: lsp.RPCVersion,
				ID:  req.ID,
			},
		}
		text, ok := documents.Get(req.Params.TextDocument.URI)
		if !ok {
			return nil, nil
		}
		urls, ok := urls.Get(req.Params.TextDocument.URI)
		if !ok {
			return nil, nil
		}
		if len(*urls) == 0 {
			return nil, nil
		}
		h, err := db.Queries.GetHTMLByURL(
			ctx,
			master.GetHTMLByURLParams{Value: (*urls)[0]},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get html: %w", err)
		}
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(h.Value))
		if err != nil {
			return nil, fmt.Errorf("failed to get html: %w", err)
		}
		response.Result, err = GetSelectorHover(
			ctx,
			req.Params.Position,
			text,
			doc,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get hover: %w", err)
		}
		return response, nil
	}
}

// GetSelectorHover checks if the position is within the struct tag.
func GetSelectorHover(
	ctx context.Context,
	position protocol.Position,
	text *string,
	doc *goquery.Document,
) (res lsp.HoverResult, err error) {
	select {
	case <-ctx.Done():
		return res, fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		var inValue bool
		fset := token.NewFileSet()
		position.Line = position.Line + 1
		node, err := parser.ParseFile(
			fset,
			"",
			bytes.NewBufferString(*text),
			parser.Trace,
		)
		if err != nil {
			return res, fmt.Errorf("failed to parse struct: %w", err)
		}
		structNodes := parsers.FindStructNodes(node)
		var resCh chan lsp.HoverResult = make(chan lsp.HoverResult)
		for i := range structNodes {
			go func(i int) {
				// Check if the position is within the struct node
				inPosition := parsers.IsPositionInNode(
					structNodes[i],
					position,
					fset,
				)
				// Check if the position is within a struct tag
				inTag := parsers.IsPositionInTag(
					structNodes[i],
					position,
					fset,
				)
				if !inPosition && !inTag {
					return
				}
				var val string
				// Check if the position is within a struct tag value
				// (i.e. value inside and including " and " characters)
				val, inValue = parsers.PositionInStructTagValue(
					structNodes[i],
					position,
					fset,
					text,
				)
				if !inValue && val != "" {
					docHTML, err := doc.Html()
					if err != nil {
						log.Errorf("failed to get html: %s", err)
					}
					val = fmt.Sprintf(
						"`%s`\n%s",
						val,
						docHTML,
					)
					res.Contents = val
					resCh <- res
					return
				}
				var HTMLs []string
				found := doc.Find(val)
				HTMLs = make([]string, found.Length())
				eg := errgroup.Group{}
				found.Each(func(i int, s *goquery.Selection) {
					eg.Go(func() error {
						HTML, err := s.Parent().Html()
						if err != nil {
							return fmt.Errorf("failed to get html: %w", err)
						}
						HTMLs[i] = fmt.Sprintf(
							"%d:\n%s",
							i,
							gohtml.Format(HTML),
						)
						return nil
					})
				})
				err := eg.Wait()
				if err != nil {
					log.Errorf("failed to get html: %s", err)
				}
				HTML := strings.Join(HTMLs, "\n================\n")
				res.Contents = fmt.Sprintf(
					"`%s`:\n%s",
					val,
					HTML,
				)
				resCh <- res
			}(i)
		}
		return <-resCh, nil
	}
}
