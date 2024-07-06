package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		if err = json.Unmarshal([]byte(contents), &request); err != nil {
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
		if err = json.Unmarshal([]byte(contents), &request); err != nil {
			return fmt.Errorf(
				"decode (initialized) request failed: %w",
				err,
			)
		}
	case methods.MethodRequestTextDocumentDidOpen:
		var request lsp.NotificationDidOpenTextDocument
		if err = json.Unmarshal(contents, &request); err != nil {
			return fmt.Errorf(
				"decode (textDocument/didOpen) request failed: %w",
				err,
			)
		}
		response, err := state.OpenDocument(
			ctx,
			request,
		)
		if err != nil {
			return fmt.Errorf("failed to open document: %w", err)
		}
		err = WriteResponse(ctx, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case methods.MethodTextDocumentDidClose:
		var request lsp.DidCloseTextDocumentParamsNotification
		if err = json.Unmarshal([]byte(contents), &request); err != nil {
			return fmt.Errorf("decode (didClose) request failed: %w", err)
		}
		state.Documents[request.Params.TextDocument.URI] = ""
	case methods.MethodRequestTextDocumentCompletion:
		var request lsp.CompletionRequest
		err = json.Unmarshal(contents, &request)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal completion request (textDocument/completion): %w",
				err,
			)
		}
		response, err := state.CreateTextDocumentCompletion(
			request.ID,
			request.Params.TextDocument,
			request.Params.Position,
		)
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
		diagnostics := []lsp.Diagnostic{}
		for _, change := range request.Params.ContentChanges {
			diags, err := state.UpdateDocument(
				request.Params.TextDocument.URI,
				change.Text,
			)
			if err != nil {
				return fmt.Errorf("failed to update document: %w", err)
			}
			diagnostics = append(diagnostics, diags...)
		}
		response := lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
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
		response, err := state.Hover(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		)
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
		response, err := state.TextDocumentCodeAction(
			request,
		)
		if err != nil {
			return fmt.Errorf("failed to get code actions: %w", err)
		}
		err = WriteResponse(ctx, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case methods.MethodTextDocumentDidSave:
		// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didSave
		state.Logger.Printf("Client sent a did save notification")
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
		// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#$/cancelRequest
		var request lsp.CancelRequest
		err = json.Unmarshal(contents, &request)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal cancel request ($/cancelRequest): %w",
				err,
			)
		}
		response, err := state.CancelRequest(request.ID)
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
	default:
		return fmt.Errorf("unknown method: %s", method)
	}
	return nil
}
