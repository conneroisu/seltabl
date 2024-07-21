package analysis

import (
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
)

// CancelRequest cancels a request
func (s *State) CancelRequest(
	request lsp.CancelRequest,
) (response *lsp.CancelResponse, err error) {
	return &lsp.CancelResponse{
		Response: lsp.Response{
			RPC: lsp.RPCVersion,
			ID:  int(request.ID),
		},
	}, nil
}
