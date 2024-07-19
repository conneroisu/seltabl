package lsp

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
)

// InitializeRequest is a struct for the initialize request.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialize
type InitializeRequest struct {
	// InitializeRequest embeds the Request struct
	Request
	// Params are the parameters for the initialize request.
	Params InitializeRequestParams `json:"params"`
}

// Method returns the method for the initialize request.
func (r InitializeRequest) Method() methods.Method {
	return methods.MethodInitialize
}

// InitializeRequestParams is a struct for the initialize request params
type InitializeRequestParams struct {
	// ClientInfo is the client info of the client in the request
	ClientInfo *ClientInfo `json:"clientInfo"`
	// InitializationOptions is the initialization options of the client in the request
	RootPath string `json:"rootPath,omitempty"`
	// Trace is the trace of the client in the request
	Trace string `json:"trace,omitempty"`
	// ProcessID is the process id of the client in the request
	ProcessID int `json:"processId,omitempty"`
	// Locale is the locale of the client in the request
	Locale string `json:"locale,omitempty"`
	// RootURI is the root uri of the client in the request
	RootURI string `json:"rootUri,omitempty"`
	// WorkspaceFolders is the workspace folders of the client in the request
	WorkspaceFolders []WorkspaceFolder `json:"workspaceFolders,omitempty"`
}

// WorkspaceFolder is a struct for the workspace folder
type WorkspaceFolder struct {
	// URI is the uri of the workspace folder
	URI string `json:"uri"`
	// Name is the name of the workspace folder
	Name string `json:"name"`
}

// ClientInfo is a struct for the client info
type ClientInfo struct {
	// Name is the name of the client
	Name string `json:"name"`
	// Version is the version of the client
	Version string `json:"version"`
}

// InitializeResponse is a struct for the initialize response.
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialize
//
// It embeds the Response struct.
type InitializeResponse struct {
	Response
	// Result is the result of the initialize request
	Result InitializeResult `json:"result"`
}

// Method returns the method for the initialize response
func (r InitializeResponse) Method() methods.Method {
	return methods.MethodInitialize
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
	TextDocumentSync   int            `json:"textDocumentSync"`   // TextDocumentSync is what the server supports for syncing text documents.
	HoverProvider      bool           `json:"hoverProvider"`      // HoverProvider is a boolean indicating whether the server provides.
	DefinitionProvider bool           `json:"definitionProvider"` // DefinitionProvider is a boolean indicating whether the server provides definition capabilities.
	CodeActionProvider bool           `json:"codeActionProvider"` // CodeActionProvider is a boolean indicating whether the server provides code actions.
	CompletionProvider map[string]any `json:"completionProvider"` // CompletionProvider is a map of completion providers.
}

// ServerInfo is a struct for the server info.
type ServerInfo struct {
	// Name is the name of the server
	Name string `json:"name"`
	// Version is the version of the server
	Version string `json:"version"`
}

// NewInitializeResponse creates a new initialize response.
func NewInitializeResponse(
	ctx context.Context,
	request *InitializeRequest,
) (*InitializeResponse, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
			return &InitializeResponse{
				Response: Response{
					RPC: RPCVersion,
					ID:  request.ID,
				},
				Result: InitializeResult{
					Capabilities: ServerCapabilities{
						TextDocumentSync:   1,
						HoverProvider:      true,
						DefinitionProvider: true,
						CodeActionProvider: false,
						CompletionProvider: map[string]any{},
					},
					ServerInfo: ServerInfo{
						Name:    "seltabl_lsp",
						Version: "0.0.0.5.0.0-beta1.final",
					},
				},
			}, nil
		}
	}
}

// InitializedParamsRequest is a struct for the initialized params.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialized
type InitializedParamsRequest struct {
	// InitializedParamsRequest embeds the Request struct
	Response
}

// Method returns the method for the initialized params request.
func (r InitializedParamsRequest) Method() methods.Method {
	return methods.MethodNotificationInitialized
}
