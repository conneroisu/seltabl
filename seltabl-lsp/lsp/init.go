package lsp

// InitializeRequest is a struct for the initialize request
type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

// InitializeRequestParams is a struct for the initialize request params
type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
}

// ClientInfo is a struct for the client info
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResponse is a struct for the initialize response
type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

// InitializeResult is a struct for the initialize result
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

// ServerCapabilities is a struct for the server capabilities
type ServerCapabilities struct {
	TextDocumentSync int `json:"textDocumentSync"`

	HoverProvider      bool           `json:"hoverProvider"`
	DefinitionProvider bool           `json:"definitionProvider"`
	CodeActionProvider bool           `json:"codeActionProvider"`
	CompletionProvider map[string]any `json:"completionProvider"`
}

// ServerInfo is a struct for the server info
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// NewInitializeResponse returns a new initialize response
func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:   1,
				HoverProvider:      true,
				DefinitionProvider: true,
				CodeActionProvider: true,
				CompletionProvider: map[string]any{},
			},
			ServerInfo: ServerInfo{
				Name:    "educationalsp",
				Version: "0.0.0.0.0.0-beta1.final",
			},
		},
	}
}
