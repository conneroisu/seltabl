package analysis

import (
	"context"
	"fmt"
	"strings"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/safe"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

// TextDocumentCodeAction returns the code actions for a given text document.
func TextDocumentCodeAction(
	ctx context.Context,
	req lsp.TextDocumentCodeActionRequest,
	documents *safe.Map[uri.URI, string],
) (response rpc.MethodActor, err error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			text, ok := documents.Get(req.Params.TextDocument.URI)
			if !ok {
				return nil, fmt.Errorf("document not found")
			}
			actions := []protocol.CodeAction{}
			for row, line := range strings.Split(*text, "\n") {
				idx := strings.Index(line, "VS Code")
				if idx >= 0 {
					replaceChange := map[string][]protocol.TextEdit{}
					replaceChange[string(req.Params.TextDocument.URI)] = []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      uint32(row),
									Character: uint32(row),
								},
								End: protocol.Position{
									Line:      uint32(row),
									Character: uint32(row),
								},
							},
							NewText: "Neovim",
						},
					}
					var a = make(map[uri.URI][]protocol.TextEdit)
					a[uri.URI(req.Params.TextDocument.URI)] = []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      uint32(row),
									Character: uint32(row),
								},
								End: protocol.Position{
									Line:      uint32(row),
									Character: uint32(row),
								},
							},
							NewText: "Neovim",
						},
					}
					actions = append(actions, protocol.CodeAction{
						Title: "Replace VS C*de with a superior editor",
						Edit:  &protocol.WorkspaceEdit{Changes: a},
					})
					censorChange := map[string][]protocol.TextEdit{}
					censorChange[string(req.Params.TextDocument.URI)] = []protocol.TextEdit{
						{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      uint32(row),
									Character: uint32(row),
								},
								End: protocol.Position{
									Line:      uint32(row),
									Character: uint32(row),
								},
							},
							NewText: "VS C*de",
						},
					}
					actions = append(actions, protocol.CodeAction{
						Title: "Censor to VS C*de",
						Kind:  protocol.QuickFix,
						Edit:  &protocol.WorkspaceEdit{Changes: a},
					})
				}
			}
			return &lsp.TextDocumentCodeActionResponse{
				Response: lsp.Response{
					RPC: lsp.RPCVersion,
					ID:  req.ID,
				},
				Result: actions,
			}, nil
		}
	}
}
