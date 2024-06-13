package lsp

// Request is the request to a LSP
type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`
}

// Response is the response of a LSP
type Response struct {
	RPC string `json:"jsonrpc"`
	ID  *int   `json:"id,omitempty"`
	// Result string `json:"result"`
	// Error  string `json:"error"`
}
