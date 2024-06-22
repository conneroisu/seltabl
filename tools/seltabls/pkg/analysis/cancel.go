package analysis

import (
	"strconv"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
)

// CancelRequest cancels a request
func (s *State) CancelRequest(
	id string,
) (response *lsp.CancelResponse, err error) {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return &lsp.CancelResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  int(i),
		},
	}, nil
}
