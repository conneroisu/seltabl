package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp/methods"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/safe"
	"go.lsp.dev/uri"
)

// HandleMessage handles a message sent from the client to the language server.
// It parses the message and returns with a response.
func HandleMessage(
	ctx context.Context,
	cancel *context.CancelFunc,
	msg *rpc.BaseMessage,
	db *data.Database[master.Queries],
	documents *safe.Map[uri.URI, string],
	selectors *safe.Map[uri.URI, []master.Selector],
	urls *safe.Map[uri.URI, []string],
) (response rpc.MethodActor, err error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		switch methods.Method(msg.Method) {
		case methods.MethodInitialize:
			var request lsp.InitializeRequest
			err = json.Unmarshal([]byte(msg.Content), &request)
			if err != nil {
				return nil, fmt.Errorf(
					"decode initialize request (initialize) failed: %w",
					err,
				)
			}
			return lsp.NewInitializeResponse(ctx, &request)
		case methods.MethodRequestTextDocumentDidOpen:
			var request lsp.NotificationDidOpenTextDocument
			err = json.Unmarshal(msg.Content, &request)
			if err != nil {
				return nil, fmt.Errorf(
					"decode (textDocument/didOpen) request failed: %w",
					err,
				)
			}
			if !strings.HasSuffix(
				string(request.Params.TextDocument.URI),
				".go",
			) {
				return nil, nil
			}
			return analysis.OpenDocument(
				ctx,
				&request,
				db,
				documents,
				urls,
				selectors,
			)
		case methods.MethodRequestTextDocumentCompletion:
			var request lsp.TextDocumentCompletionRequest
			err = json.Unmarshal(msg.Content, &request)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal completion request (textDocument/completion): %w",
					err,
				)
			}
			return analysis.CreateTextDocumentCompletion(
				ctx,
				request,
				documents,
				selectors,
			)
		case methods.MethodRequestTextDocumentHover:
			var request lsp.HoverRequest
			err = json.Unmarshal(msg.Content, &request)
			if err != nil {
				return nil, fmt.Errorf(
					"failed unmarshal of hover request (): %w",
					err,
				)
			}
			return analysis.NewHoverResponse(
				ctx,
				request,
				documents,
				urls,
			)
		// case methods.MethodRequestTextDocumentCodeAction:
		//         var request lsp.CodeActionRequest
		//         err = json.Unmarshal(msg.Content, &request)
		//         if err != nil {
		//                 return nil, fmt.Errorf(
		//                         "failed to unmarshal of codeAction request (textDocument/codeAction): %w",
		//                         err,
		//                 )
		//         }
		//         return analysis.TextDocumentCodeAction(
		//                 ctx,
		//                 request,
		//                 state,
		//         )
		case methods.MethodShutdown:
			var request lsp.ShutdownRequest
			err = json.Unmarshal([]byte(msg.Content), &request)
			if err != nil {
				return nil, fmt.Errorf(
					"decode (shutdown) request failed: %w",
					err,
				)
			}
			(*cancel)()
			return lsp.NewShutdownResponse(request, nil)
		case methods.MethodCancelRequest:
			var request lsp.CancelRequest
			err = json.Unmarshal(msg.Content, &request)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to unmarshal cancel request ($/cancelRequest): %w",
					err,
				)
			}
			log.Debugf("canceling request: %d", request.Params.ID)
			c, ok := lsp.CancelMap.Get(int(request.Params.ID.(float64)))
			if ok {
				(*c)()
			}
			return analysis.CancelResponse(request)
		case methods.MethodNotificationInitialized:
			var request lsp.InitializedParamsRequest
			err = json.Unmarshal([]byte(msg.Content), &request)
			if err != nil {
				return nil, fmt.Errorf(
					"decode (initialized) request failed: %w",
					err,
				)
			}
			return nil, nil
		case methods.MethodNotificationExit:
			os.Exit(0)
			return nil, nil
		case methods.MethodNotificationTextDocumentWillSave:
			var request lsp.WillSaveTextDocumentNotification
			err = json.Unmarshal([]byte(msg.Content), &request)
			if err != nil {
				return nil, fmt.Errorf(
					"decode (willSave) request failed: %w",
					err,
				)
			}
			return nil, nil
		case methods.MethodNotificationTextDocumentDidSave:
			var request lsp.DidSaveTextDocumentNotification
			err = json.Unmarshal([]byte(msg.Content), &request)
			if err != nil {
				return nil, fmt.Errorf(
					"decode (didSave) request failed: %w",
					err,
				)
			}
			read, err := ReadFile(request.Params.TextDocument.URI)
			if err != nil {
				return nil, fmt.Errorf("failed to read file: %w", err)
			}
			documents.Set(request.Params.TextDocument.URI, string(read))
			return nil, nil
		case methods.NotificationTextDocumentDidClose:
			var request lsp.DidCloseTextDocumentParamsNotification
			if err = json.Unmarshal([]byte(msg.Content), &request); err != nil {
				return nil, fmt.Errorf(
					"decode (didClose) request failed: %w",
					err,
				)
			}
			return nil, nil
		case methods.NotificationMethodTextDocumentDidChange:
			var request lsp.TextDocumentDidChangeNotification
			err = json.Unmarshal(msg.Content, &request)
			if err != nil {
				return nil, fmt.Errorf(
					"decode (textDocument/didChange) request failed: %w",
					err,
				)
			}
			return analysis.UpdateDocument(
				ctx,
				&request,
				db,
				documents,
				urls,
				selectors,
			)
		default:
			return nil, fmt.Errorf("unknown method: %s", msg.Method)
		}
	}
}

// ReadFile reads a file from the given uri.
func ReadFile(uri uri.URI) (string, error) {
	u, err := url.Parse(string(uri))
	if err != nil {
		return "", fmt.Errorf("failed to parse uri: %w", err)
	}
	read, err := os.ReadFile(u.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(read), nil
}
