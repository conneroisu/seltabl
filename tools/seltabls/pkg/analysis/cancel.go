package analysis

import (
	"strconv"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
)

// CancelRequest cancels a request
func (s *State) CancelRequest(
	request lsp.CancelRequest,
) (response *lsp.CancelResponse, err error) {
	i, err := strconv.ParseInt(request.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	return &lsp.CancelResponse{
		Response: lsp.Response{
			RPC: lsp.RPCVersion,
			ID:  int(i),
		},
	}, nil
}
