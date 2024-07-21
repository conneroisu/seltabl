package analysis

import (
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc/errors"
)

// CancelRequest cancels a request
func (s *State) CancelRequest(
	request lsp.CancelRequest,
) (response *lsp.CancelResponse, err error) {
	return &lsp.CancelResponse{
		RPC: lsp.RPCVersion,
		ID:  request.Params.ID,
		Error: &lsp.Error{
			Code:    int(errors.CodeRequestCancelled),
			Message: "Request cancelled",
		},
	}, nil
}
