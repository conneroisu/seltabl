package lsp

// InitializeRequest is a struct for the initialize request.
type InitializeRequest struct {
	// InitializeRequest embeds the Request struct
	Request
	// Params are the parameters for the initialize request.
	Params InitializeRequestParams `json:"params"`
}

// InitializeRequestParams is a struct for the initialize request params
type InitializeRequestParams struct {
	// ClientInfo is the client info of the client in the request
	ClientInfo *ClientInfo `json:"clientInfo"`
	// InitializationOptions is the initialization options of the client in the request
	RootPath string `json:"rootPath,omitempty"`
	// Trace is the trace of the client in the request
	Trace string `json:"trace,omitempty"`
}

// ClientInfo is a struct for the client info
type ClientInfo struct {
	// Name is the name of the client
	Name string `json:"name"`
	// Version is the version of the client
	Version string `json:"version"`
}

// InitializeResponse is a struct for the initialize response.
// It embeds the Response struct.
type InitializeResponse struct {
	Response
	// Result is the result of the initialize request
	Result InitializeResult `json:"result"`
}

// InitializeResult is a struct for the initialize result used in the initialize response.
type InitializeResult struct {
	// Capabilities are the capabilities of the server for the initialize response.
	Capabilities ServerCapabilities `json:"capabilities"`
	// ServerInfo is the server info for the initialize response.
	ServerInfo ServerInfo `json:"serverInfo"`
}

// ServerCapabilities is a struct for the server capabilities
type ServerCapabilities struct {
	// TextDocumentSync is what the server supports for syncing text documents.
	TextDocumentSync int `json:"textDocumentSync"`
	// HoverProvider is a boolean indicating whether the server provides.
	HoverProvider bool `json:"hoverProvider"`
	// DefinitionProvider is a boolean indicating whether the server provides definition capabilities.
	DefinitionProvider bool `json:"definitionProvider"`
	// CodeActionProvider is a boolean indicating whether the server provides code actions.
	CodeActionProvider bool `json:"codeActionProvider"`
	// CompletionProvider is a map of completion providers.
	CompletionProvider map[string]any `json:"completionProvider"`
}

// ServerInfo is a struct for the server info.
type ServerInfo struct {
	// Name is the name of the server
	Name string `json:"name"`
	// Version is the version of the server
	Version string `json:"version"`
}

// NewInitializeResponse creates a new initialize response.
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
				Name:    "seltabl_lsp",
				Version: "0.0.0.0.0.0-beta1.final",
			},
		},
	}
}

// InitializedParamsRequest is a struct for the initialized params.
type InitializedParamsRequest struct {
	// InitializedParamsRequest embeds the Request struct
	Response
}

// InitializedParamsResponse is a struct for the initialized params.
type InitializedParamsResponse struct {
	// InitializedParamsResponse embeds the Response struct
	Response
	// Params are the parameters for the initialized params.
	Params InitializedParams `json:"params"`
}

// InitializedParams is a struct for the initialized params.
type InitializedParams struct {
	// The process ID of the client for which the server has been started.
	ProcessID int `json:"processId"`
	// The root path of the workspace. This is the path of the folder that contains the workspace file.
	RootPath string `json:"rootPath,omitempty"`
	// The root URI of the workspace. This is the URI of the folder that contains the workspace file.
	RootURI string `json:"rootUri,omitempty"`
	// The initialization options of the client in the request.
	Trace string `json:"trace,omitempty"`
	// The capabilities provided by the client in the request.
	Capabilities ClientCapabilities `json:"capabilities"`
}

// NewInitializedParamsResponse returns a new initialized params response
func NewInitializedParamsResponse(id int) InitializedParamsResponse {
	return InitializedParamsResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Params: InitializedParams{
			// The process ID of the client for which the server has been started.
			ProcessID: 1,
			// The initialization options of the client in the request.
			Trace: "off",
			// The capabilities provided by the client in the request.
			Capabilities: ClientCapabilities{
				TextDocumentSync:   1,
				HoverProvider:      true,
				DefinitionProvider: true,
				CodeActionProvider: true,
				CompletionProvider: map[string]any{},
			},
		},
	}
}
