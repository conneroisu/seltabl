package analysis

import (
	"context"
	"fmt"
	"strings"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

// TextDocumentCodeAction returns the code actions for a given text document.
func TextDocumentCodeAction(
	ctx context.Context,
	req lsp.CodeActionRequest,
	s *State,
) (response rpc.MethodActor, err error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			text := s.Documents[string(req.Params.TextDocument.URI)]
			actions := []protocol.CodeAction{}
			for row, line := range strings.Split(text, "\n") {
				idx := strings.Index(line, "VS Code")
				if idx >= 0 {
					replaceChange := map[string][]lsp.TextEdit{}
					replaceChange[string(req.Params.TextDocument.URI)] = []lsp.TextEdit{
						{
							Range: lsp.LineRange(
								row,
								idx,
								idx+len("VS Code"),
							),
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
					censorChange := map[string][]lsp.TextEdit{}
					censorChange[string(req.Params.TextDocument.URI)] = []lsp.TextEdit{
						{
							Range: lsp.LineRange(
								row,
								idx,
								idx+len("VS Code"),
							),
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
