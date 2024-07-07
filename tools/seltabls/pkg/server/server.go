package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
)

// HandleMessage handles a message sent from the client to the language server.
// It parses the message and returns with a response.
func HandleMessage(
	ctx context.Context,
	writer *io.Writer,
	state *analysis.State,
	msg rpc.BaseMessage,
) (err error) {
	method := msg.Method
	contents := msg.Content
	switch methods.GetMethod(method) {
	case methods.MethodInitialize:
		var request lsp.InitializeRequest
		err = json.Unmarshal([]byte(contents), &request)
		if err != nil {
			return fmt.Errorf(
				"decode initialize request (initialize) failed: %w",
				err,
			)
		}
		response := lsp.NewInitializeResponse(&request)
		err = WriteResponse(ctx, writer, response)
		if err != nil {
			return fmt.Errorf(
				"failed to write (initialize) response: %w",
				err,
			)
		}
	case methods.MethodNotificationInitialized:
		var request lsp.InitializedParamsRequest
		err = json.Unmarshal([]byte(contents), &request)
		if err != nil {
			return fmt.Errorf(
				"decode (initialized) request failed: %w",
				err,
			)
		}
	case methods.MethodRequestTextDocumentDidOpen:
		var request lsp.NotificationDidOpenTextDocument
		err = json.Unmarshal(contents, &request)
		if err != nil {
			return fmt.Errorf(
				"decode (textDocument/didOpen) request failed: %w",
				err,
			)
		}
		response, err := state.OpenDocument(ctx, request)
		if err != nil {
			return fmt.Errorf("failed to open document: %w", err)
		}
		if len(response.Params.Diagnostics) == 0 {
			return nil
		}
		err = WriteResponse(ctx, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case methods.MethodRequestTextDocumentCompletion:
		var request lsp.CompletionRequest
		err = json.Unmarshal(contents, &request)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal completion request (textDocument/completion): %w",
				err,
			)
		}
		response, err := state.CreateTextDocumentCompletion(ctx, request)
		if err != nil {
			return fmt.Errorf("failed to get completions: %w", err)
		}
		err = WriteResponse(ctx, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case methods.NotificationMethodTextDocumentDidChange:
		var request lsp.TextDocumentDidChangeNotification
		err = json.Unmarshal(contents, &request)
		if err != nil {
			return fmt.Errorf(
				"decode (textDocument/didChange) request failed: %w",
				err,
			)
		}
		response, err := analysis.UpdateDocument(ctx, state, &request)
		if err != nil {
			return fmt.Errorf("failed to update document: %w", err)
		}
		err = WriteResponse(ctx, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case methods.MethodRequestTextDocumentHover:
		var request lsp.HoverRequest
		err = json.Unmarshal(contents, &request)
		if err != nil {
			return fmt.Errorf("failed unmarshal of hover request (): %w", err)
		}
		response, err := analysis.NewHoverResponse(request, state)
		if err != nil {
			return fmt.Errorf("failed to get hover: %w", err)
		}
		err = WriteResponse(ctx, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case methods.MethodRequestTextDocumentCodeAction:
		var request lsp.CodeActionRequest
		err = json.Unmarshal(contents, &request)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal of codeAction request (textDocument/codeAction): %w",
				err,
			)
		}
		response, err := analysis.TextDocumentCodeAction(request, state)
		if err != nil {
			return fmt.Errorf("failed to get code actions: %w", err)
		}
		err = WriteResponse(ctx, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case methods.MethodShutdown:
		var request lsp.ShutdownRequest
		err = json.Unmarshal([]byte(contents), &request)
		if err != nil {
			return fmt.Errorf("decode (shutdown) request failed: %w", err)
		}
		response := lsp.NewShutdownResponse(request, nil)
		err = WriteResponse(ctx, writer, response)
		if err != nil {
			return fmt.Errorf("write (shutdown) response failed: %w", err)
		}
		os.Exit(0)
	case methods.MethodCancelRequest:
		var request lsp.CancelRequest
		err = json.Unmarshal(contents, &request)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal cancel request ($/cancelRequest): %w",
				err,
			)
		}
		response, err := state.CancelRequest(request)
		if err != nil {
			return fmt.Errorf("failed to cancel request: %w", err)
		}
		err = WriteResponse(ctx, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case methods.MethodExit:
		os.Exit(0)
		return nil
	case methods.MethodNotificationTextDocumentDidSave:
		state.Logger.Printf("Client sent a did save notification")
		var request lsp.DidSaveTextDocumentParamsNotification
		err = json.Unmarshal([]byte(contents), &request)
		if err != nil {
			return fmt.Errorf("decode (didSave) request failed: %w", err)
		}
		u, err := url.Parse(request.Params.TextDocument.URI)
		content, err := os.ReadFile(u.Path)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		// repsond with diagnostics for the file
		diagsResp, err := state.OpenDocument(ctx, lsp.NotificationDidOpenTextDocument{
			Notification: lsp.Notification{
				RPC:    lsp.RPCVersion,
				Method: "textDocument/didOpen",
			},
			Params: lsp.DidOpenTextDocumentParams{
				TextDocument: lsp.TextDocumentItem{
					URI:        request.Params.TextDocument.URI,
					Text:       string(content),
					LanguageID: "go",
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to get diagnostics for file: %w", err)
		}
		err = WriteResponse(ctx, writer, diagsResp)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
		state.Logger.Printf("Client completed a did save notification")
	case methods.MethodNotificationTextDocumentDidClose:
		var request lsp.DidCloseTextDocumentParamsNotification
		if err = json.Unmarshal([]byte(contents), &request); err != nil {
			return fmt.Errorf("decode (didClose) request failed: %w", err)
		}
		state.Documents[request.Params.TextDocument.URI] = ""
	default:
		return fmt.Errorf("unknown method: %s", method)
	}
	return nil
}
