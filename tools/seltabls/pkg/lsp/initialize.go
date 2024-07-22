package lsp

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"go.lsp.dev/protocol"
)

// InitializeRequest is a struct for the initialize request.
//
// Microsoft LSP Docs:
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#initialize
type InitializeRequest struct {
	// InitializeRequest embeds the Request struct
	Request
	// Params are the parameters for the initialize request.
	Params protocol.InitializeParams `json:"params"`
}

// Method returns the method for the initialize request.
func (r InitializeRequest) Method() methods.Method {
	return methods.MethodInitialize
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
	Result protocol.InitializeResult `json:"result"`
}

// Method returns the method for the initialize response
func (r InitializeResponse) Method() methods.Method {
	return methods.MethodInitialize
}

// General is a struct for the general capabilities.
type General struct {
	SupportsCancellation bool `json:"supportsCancellation"`
}

// ServerCapabilities is a struct for the server capabilities
type ServerCapabilities struct {
	TextDocumentSync      int                   `json:"textDocumentSync"`                // TextDocumentSync is what the server supports for syncing text documents.
	HoverProvider         bool                  `json:"hoverProvider"`                   // HoverProvider is a boolean indicating whether the server provides.
	DefinitionProvider    bool                  `json:"definitionProvider"`              // DefinitionProvider is a boolean indicating whether the server provides definition capabilities.
	CodeActionProvider    bool                  `json:"codeActionProvider"`              // CodeActionProvider is a boolean indicating whether the server provides code actions.
	CompletionProvider    map[string]any        `json:"completionProvider"`              // CompletionProvider is a map of completion providers.
	SignatureHelpProvider *SignatureHelpOptions `json:"signatureHelpProvider,omitempty"` // SignatureHelpProvider is a boolean indicating whether the server provides signature help.
	CancellationProvider  bool                  `json:"cancellationProvider,omitempty"`  // CancellationProvider is a boolean indicating whether the server supports cancellation.
}

type SignatureHelpOptions struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
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
				Result: protocol.InitializeResult{
					Capabilities: protocol.ServerCapabilities{
						TextDocumentSync: protocol.TextDocumentSyncOptions{
							OpenClose: true,
							Change:    protocol.TextDocumentSyncKindFull,
						},
						CompletionProvider: &protocol.CompletionOptions{
							ResolveProvider: false,
							TriggerCharacters: []string{
								":", "\"",
							},
						},
						HoverProvider:      true,
						DefinitionProvider: false,
						CodeActionProvider: false,
					},
					ServerInfo: &protocol.ServerInfo{
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
