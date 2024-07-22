package analysis

import (
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
)

// CancelRequest cancels a request
func (s *State) CancelRequest(
	request lsp.CancelRequest,
) (response *lsp.CancelResponse, err error) {
	return &lsp.CancelResponse{
		RPC: lsp.RPCVersion,
		ID:  request.Params.ID,
	}, nil
}
