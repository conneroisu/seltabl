package analysis

import (
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
)

// CancelResponse cancels a request
func CancelResponse(
	request lsp.CancelRequest,
) (response *lsp.CancelResponse, err error) {
	return &lsp.CancelResponse{
		RPC: lsp.RPCVersion,
		ID:  int(request.Params.ID.(float64)),
	}, nil
}
