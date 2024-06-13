package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/conneroisu/seltabl/tools/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/pkg/rpc"
)

// handleMessage handles a message sent from the client to the language server.
// It parses the message and returns with a response.
func handleMessage(
	logger *log.Logger,
	writer io.Writer,
	state *analysis.State,
	method string,
	contents []byte,
) error {
	var err error
	logger.Printf("Received message (%s): %s\n", method, contents)
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("Recovered from panic: %v\n", r)
		}
	}()
	defer func() {
		if err != nil {
			logger.Printf("Error: %v\n", err)
		}
	}()
	switch method {
	case "initialize":
		logger.Println("received initialize request: ", contents)
		var request lsp.InitializeRequest
		if err = json.Unmarshal([]byte(contents), &request); err != nil {
			logger.Printf("failed to decode initialize request (initialize): %s\n", err)
			return fmt.Errorf("failed to decode initialize request (initialize): %w", err)
		}
		logger.Printf(
			"Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version,
		)
		msg := lsp.NewInitializeResponse(request.ID)
		if err = writeResponse(writer, msg); err != nil {
			logger.Print(fmt.Errorf("failed to write a response: %v\nthe response that failed to write: %s\n", err, contents))
			return fmt.Errorf("failed to write response: %w", err)
		}
		enc, err := rpc.EncodeMessage(msg)
		if err != nil {
			logger.Printf("failed to encode message: %s\n", err)
			return fmt.Errorf("failed to encode message: %w", err)
		}
		logger.Println("Received message (initialize) and replied: ", enc)
	case "initialized":
		logger.Println("received initialized request: ", contents)
		var request lsp.InitializedParamsRequest
		if err = json.Unmarshal([]byte(contents), &request); err != nil {
			logger.Printf("failed to decode initialized request (initialized): %s\n", err)
			return fmt.Errorf("failed to decode initialized request (initialized): %w", err)
		}
		msg := lsp.NewInitializedParamsResponse(*request.ID)
		if err = writeResponse(writer, msg); err != nil {
			logger.Print(fmt.Errorf("failed to write a response: %v\nthe response that failed to write: %s\n", err, contents))
			return fmt.Errorf("failed to write response: %w", err)
		}
		enc, err := rpc.EncodeMessage(msg)
		if err != nil {
			logger.Printf("failed to encode message: %s\n", err)
			return fmt.Errorf("failed to encode message: %w", err)
		}
		logger.Println("Received message (initialized) and replied: ", enc)
	case "textDocument/didClose":
		logger.Println("Received message: Content-Length: 3587\n", contents)
		var request lsp.DidCloseTextDocumentParamsNotification
		if err = json.Unmarshal([]byte(contents), &request); err != nil {
			logger.Printf("failed to decode didClose request (didClose): %s\n", err)
			return fmt.Errorf("failed to decode didClose request (didClose): %w", err)
		}
		msg := lsp.NewDidCloseTextDocumentParamsNotification()
		if err = writeResponse(writer, msg); err != nil {
			logger.Print(fmt.Errorf("failed to write a response: %v\nthe response that failed to write: %s\n", err, contents))
			return fmt.Errorf("failed to write response: %w", err)
		}
		enc, err := rpc.EncodeMessage(msg)
		if err != nil {
			logger.Printf("failed to encode message: %s\n", err)
			return fmt.Errorf("failed to encode message: %w", err)
		}
		logger.Println("Received message (didClose) and replied: ", enc)
	case "textDocument/didOpen":
		logger.Println("Received didOpen message: ", contents)
		var request lsp.DidOpenTextDocumentNotification
		if err = json.Unmarshal(contents, &request); err != nil {
			logger.Printf("failed to decode didOpen request (textDocument/didOpen): %s\n", err)
			return fmt.Errorf("failed to decode didOpen request (textDocument/didOpen): %w", err)
		}
		logger.Printf("Opened: %s", request.Params.TextDocument.URI)
		diagnostics := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
		alert := lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		}
		if err = writeResponse(writer, alert); err != nil {
			logger.Printf("failed to write a response: %s\nthe response that failed to write: %s\n", err, contents)
			return fmt.Errorf("failed to write response: %w", err)
		}
		logger.Println("Received message (textDocument/didOpen) and replied: ", alert)
	case "textDocument/didChange":
		logger.Println("Received didChange message: ", contents)
		var request lsp.TextDocumentDidChangeNotification
		if err = json.Unmarshal(contents, &request); err != nil {
			logger.Printf("failed to decode didChange request: %s\n", err)
			return fmt.Errorf("failed to decode didChange request: %w", err)
		}

		logger.Printf("Changed: %s", request.Params.TextDocument.URI)
		for _, change := range request.Params.ContentChanges {
			diagnostics := state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
			alert := lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         request.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			}
			if err = writeResponse(writer, alert); err != nil {
				logger.Printf("failed to write a response: %s\nthe response that failed to write: %s\n", err, contents)
				return fmt.Errorf("failed to write response: %w", err)
			}
			logger.Println("Received message (textDocument/didChange) and replied: ", alert)
		}
	case "textDocument/hover":
		logger.Println("Received hover message: ", contents)
		var request lsp.HoverRequest
		if err = json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/hover: %s", err)
			return fmt.Errorf("failed unmarshal of hover request (textDocument/hover): %w", err)
		}
		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		if err = writeResponse(writer, response); err != nil {
			logger.Printf("failed to write a response: %s\nthe response that failed to write: %s\n", err, contents)
			return fmt.Errorf("failed to write response: %w", err)
		}
		logger.Println("Received message (textDocument/didChange) and replied: ", contents)
	case "textDocument/definition":
		logger.Println("Received definition message: ", contents)
		var request lsp.DefinitionRequest
		if err = json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/definition: %s", err)
			return fmt.Errorf("failed unmarshal of definition request (textDocument/definition): %w", err)
		}
		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		if err = writeResponse(writer, response); err != nil {
			logger.Printf("failed to write a response: %s\nthe response that failed to write: %s\n", err, contents)
			return fmt.Errorf("failed to write response: %w", err)
		}
		logger.Println("Received message (textDocument/definition) and replied: ", response)
	case "textDocument/codeAction":
		logger.Println("Received codeAction message: ", contents)
		var request lsp.CodeActionRequest
		if err = json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/codeAction: %s", err)
			return fmt.Errorf("failed unmarshal of codeAction request (textDocument/codeAction): %w", err)
		}
		response := state.TextDocumentCodeAction(request.ID, request.Params.TextDocument.URI)
		if err = writeResponse(writer, response); err != nil {
			logger.Printf("failed to write a response: %s\nthe response that failed to write: %s\n", err, contents)
			return fmt.Errorf("failed to write response: %w", err)
		}
		logger.Println("Received message (textDocument/codeAction) and replied: ", response)
	case "textDocument/completion":
		logger.Println("Received completion message: ", contents)
		var request lsp.CompletionRequest
		if err = json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/codeAction: %s", err)
			return fmt.Errorf("failed unmarshal of completion request (textDocument/completion): %w", err)
		}
		response := state.TextDocumentCompletion(request.ID, &request.Params.TextDocument, &request.Params.Position)
		if err = writeResponse(writer, response); err != nil {
			logger.Printf("failed to write a response: %s\nthe response that failed to write: %s\n", err, contents)
			return fmt.Errorf("failed to write response: %w", err)
		}
		logger.Println("Received message (textDocument/completion) and replied: ", response)
	default:
		logger.Println("Unknown method: ", method)
		return nil
	}

	return nil
}
