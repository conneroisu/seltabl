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
	hCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case <-hCtx.Done():
			var request interface{}
			err = json.Unmarshal(msg.Content, request)
			if err != nil {
				return fmt.Errorf("failed to unmarshal cancel request ($/cancelRequest): %w", err)
			}
			cnReq, ok := request.(lsp.CancelRequest)
			if !ok {
				return fmt.Errorf("failed to cast cancel request ($/cancelRequest): %w", err)
			}
			response, err := state.CancelRequest(cnReq)
			if err != nil || response == nil {
				return fmt.Errorf("failed to cancel request: %w", err)
			}
			err = WriteResponse(hCtx, writer, response)
			if err != nil {
				return fmt.Errorf("failed to write response: %w", err)
			}
			return fmt.Errorf("context cancelled: %w", hCtx.Err())
		default:
			switch methods.GetMethod(msg.Method) {
			case methods.MethodInitialize:
				var request lsp.InitializeRequest
				err = json.Unmarshal([]byte(msg.Content), &request)
				if err != nil {
					return fmt.Errorf(
						"decode initialize request (initialize) failed: %w",
						err,
					)
				}
				response := lsp.NewInitializeResponse(&request)
				err = WriteResponse(hCtx, writer, response)
				if err != nil {
					return fmt.Errorf(
						"failed to write (initialize) response: %w",
						err,
					)
				}
			case methods.MethodNotificationInitialized:
				var request lsp.InitializedParamsRequest
				err = json.Unmarshal([]byte(msg.Content), &request)
				if err != nil {
					return fmt.Errorf(
						"decode (initialized) request failed: %w",
						err,
					)
				}
			case methods.MethodRequestTextDocumentDidOpen:
				var request lsp.NotificationDidOpenTextDocument
				err = json.Unmarshal(msg.Content, &request)
				if err != nil {
					return fmt.Errorf(
						"decode (textDocument/didOpen) request failed: %w",
						err,
					)
				}
				response, err := analysis.OpenDocument(hCtx, state, request)
				if err != nil || response == nil {
					return fmt.Errorf("failed to open document: %w", err)
				}
				err = WriteResponse(hCtx, writer, response)
				if err != nil {
					return fmt.Errorf("failed to write response: %w", err)
				}
			case methods.MethodRequestTextDocumentCompletion:
				var request lsp.CompletionRequest
				err = json.Unmarshal(msg.Content, &request)
				if err != nil {
					return fmt.Errorf(
						"failed to unmarshal completion request (textDocument/completion): %w",
						err,
					)
				}
				response, err := analysis.CreateTextDocumentCompletion(hCtx, state, request)
				if err != nil || response == nil {
					return fmt.Errorf("failed to get completions: %w", err)
				}
				err = WriteResponse(hCtx, writer, response)
				if err != nil {
					return fmt.Errorf("failed to write response: %w", err)
				}
			case methods.NotificationMethodTextDocumentDidChange:
				var request lsp.TextDocumentDidChangeNotification
				err = json.Unmarshal(msg.Content, &request)
				if err != nil {
					return fmt.Errorf(
						"decode (textDocument/didChange) request failed: %w",
						err,
					)
				}
				response, err := analysis.UpdateDocument(hCtx, state, &request)
				if err != nil || response == nil {
					return fmt.Errorf("failed to update document: %w", err)
				}
				err = WriteResponse(hCtx, writer, response)
				if err != nil {
					return fmt.Errorf("failed to write response: %w", err)
				}
			case methods.MethodRequestTextDocumentHover:
				var request lsp.HoverRequest
				err = json.Unmarshal(msg.Content, &request)
				if err != nil {
					return fmt.Errorf("failed unmarshal of hover request (): %w", err)
				}
				response, err := analysis.NewHoverResponse(request, state)
				if err != nil {
					return fmt.Errorf("failed to get hover: %w", err)
				}
				err = WriteResponse(hCtx, writer, response)
				if err != nil {
					return fmt.Errorf("failed to write response: %w", err)
				}
			case methods.MethodRequestTextDocumentCodeAction:
				var request lsp.CodeActionRequest
				err = json.Unmarshal(msg.Content, &request)
				if err != nil {
					return fmt.Errorf(
						"failed to unmarshal of codeAction request (textDocument/codeAction): %w",
						err,
					)
				}
				response, err := analysis.TextDocumentCodeAction(hCtx, request, state)
				if err != nil || response == nil {
					return fmt.Errorf("failed to get code actions: %w", err)
				}
				err = WriteResponse(hCtx, writer, response)
				if err != nil {
					return fmt.Errorf("failed to write response: %w", err)
				}
			case methods.MethodShutdown:
				var request lsp.ShutdownRequest
				err = json.Unmarshal([]byte(msg.Content), &request)
				if err != nil {
					return fmt.Errorf("decode (shutdown) request failed: %w", err)
				}
				response := lsp.NewShutdownResponse(request, nil)
				err = WriteResponse(hCtx, writer, response)
				if err != nil {
					return fmt.Errorf("write (shutdown) response failed: %w", err)
				}
				os.Exit(0)
			case methods.MethodCancelRequest:
				cancel()
				return nil
			case methods.MethodExit:
				os.Exit(0)
				return nil
			case methods.MethodNotificationTextDocumentDidSave:
				state.Logger.Printf("Client sent a did save notification")
				var request lsp.DidSaveTextDocumentParamsNotification
				err = json.Unmarshal([]byte(msg.Content), &request)
				if err != nil {
					return fmt.Errorf("decode (didSave) request failed: %w", err)
				}
				u, err := url.Parse(request.Params.TextDocument.URI)
				if err != nil {
					return fmt.Errorf("failed to parse uri: %w", err)
				}
				read, err := os.ReadFile(u.Path)
				if err != nil {
					return fmt.Errorf("failed to read file: %w", err)
				}
				state.Documents[request.Params.TextDocument.URI] = string(read)
			case methods.MethodNotificationTextDocumentDidClose:
				var request lsp.DidCloseTextDocumentParamsNotification
				if err = json.Unmarshal([]byte(msg.Content), &request); err != nil {
					return fmt.Errorf("decode (didClose) request failed: %w", err)
				}
			default:
				return fmt.Errorf("unknown method: %s", msg.Method)
			}
			return nil
		}
		return nil
	}
}
