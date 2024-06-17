package cmd

import (
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/conneroisu/seltabl/tools/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/pkg/rpc"
)

// HandleMessage handles a message sent from the client to the language server.
// It parses the message and returns with a response.
func (s *Root) HandleMessage(
	msg []byte,
) error {
	var response interface{}
	var err error
	method, contents, err := rpc.DecodeMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to decode message: %w", err)
	}
	s.Logger.Printf("Received message (%s): %s\n", method, contents)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err = sonic.Unmarshal([]byte(contents), &request); err != nil {
			return fmt.Errorf("decode initialize request (initialize) failed: %w", err)
		}
		response = lsp.NewInitializeResponse(request.ID)
		err = s.writeResponse(response)
		if err != nil {
			return fmt.Errorf("failed to write (initialize) response: %w", err)
		}
	case "initialized":
		var request lsp.InitializedParamsRequest
		if err = sonic.Unmarshal([]byte(contents), &request); err != nil {
			return fmt.Errorf("decode (initialized) request failed: %w", err)
		}
		response = lsp.NewInitializedParamsResponse(*request.ID)
		err = s.writeResponse(response)
		if err != nil {
			return fmt.Errorf("write (initialized) response failed: %w", err)
		}
	case "textDocument/didClose":
		var request lsp.DidCloseTextDocumentParamsNotification
		if err = sonic.Unmarshal([]byte(contents), &request); err != nil {
			return fmt.Errorf("decode (didClose) request failed: %w", err)
		}
		response := lsp.NewDidCloseTextDocumentParamsNotification()
		err = s.writeResponse(response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err = sonic.Unmarshal(contents, &request); err != nil {
			return fmt.Errorf("decode (textDocument/didOpen) request failed: %w", err)
		}
		diagnostics := s.State.OpenDocument(
			request.Params.TextDocument.URI,
			request.Params.TextDocument.Text,
		)
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
		err = s.writeResponse(response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		err = sonic.Unmarshal(contents, &request)
		if err != nil {
			return fmt.Errorf("decode (textDocument/didChange) request failed: %w", err)
		}
		diagnostics := []lsp.Diagnostic{}
		for _, change := range request.Params.ContentChanges {
			diagnostics = append(diagnostics, s.State.UpdateDocument(
				request.Params.TextDocument.URI,
				change.Text,
			)...)
		}
		response = lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		}
		if err = s.writeResponse(response); err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "textDocument/hover":
		var request lsp.HoverRequest
		err = sonic.Unmarshal(contents, &request)
		if err != nil {
			return fmt.Errorf("failed to unmarshal of hover request (): %w", err)
		}
		response = s.State.Hover(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		)
		err = s.writeResponse(response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "textDocument/definition":
		var request lsp.DefinitionRequest
		err = sonic.Unmarshal(contents, &request)
		if err != nil {
			return fmt.Errorf("failed to unmarshal of definition request (textDocument/definition): %w", err)
		}
		response = s.State.Definition(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		)
		err = s.writeResponse(response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		err = sonic.Unmarshal(contents, &request)
		if err != nil {
			return fmt.Errorf("failed to unmarshal of codeAction request (textDocument/codeAction): %w", err)
		}
		response = s.State.TextDocumentCodeAction(
			request.ID,
			request.Params.TextDocument.URI,
		)
		err = s.writeResponse(response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "textDocument/completion":
		var request lsp.CompletionRequest
		err = sonic.Unmarshal(contents, &request)
		if err != nil {
			return fmt.Errorf("failed to unmarshal completion request (textDocument/completion): %w", err)
		}
		response = s.State.TextDocumentCompletion(
			request.ID,
			&request.Params.TextDocument,
			&request.Params.Position,
		)
		err = s.writeResponse(response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	default:
		return fmt.Errorf("unknown method: %s", method)
	}
	enc, err := rpc.EncodeMessage(response)
	s.Logger.Println(
		fmt.Sprintf(
			"Received message (%s) err: [%w] response: `%s` contents: %s", method, err, enc, contents,
		))
	return nil
}
