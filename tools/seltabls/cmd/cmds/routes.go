package cmds

import (
	"context"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
)

func addRoutes(
	ctx context.Context,
	r *lsp.ServeMux,
	state *analysis.State,
) {
	r.Handle("initialize", initialize())
}

// initialize returns the initialize request handler.
func initialize() lsp.HandlerFunc {
	print("initialize")
	return func(w lsp.ResponseWriter, r *lsp.Request) {
		print("initialize")
		print(r.ID)
		print(r.Method)
	}
}
