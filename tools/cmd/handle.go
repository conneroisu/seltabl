package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/conneroisu/seltabl/tools/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/pkg/rpc"
)

// HandleMessage handles a message sent from the client to the language server.
// It parses the message and returns with a response.
func (s *Root) HandleMessage(
	msg []byte,
) error {
	var err error
	method, contents, err := rpc.DecodeMessage(msg)
	if err != nil {
		s.Logger.Printf("failed to decode message: %s\n", err)
		return err
	}
	s.Logger.Printf("Received message (%s): %s\n", method, contents)
	defer func() {
		if r := recover(); r != nil {
			s.Logger.Printf(
				"Recovered from panic: %v\n had error: %s\n",
				r,
				err,
			)
		}
	}()
	defer func() {
		if err != nil {
			s.Logger.Printf("Error: %v\n", err)
		}
	}()
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err = json.Unmarshal([]byte(contents), &request); err != nil {
			s.Logger.Printf(
				"failed to decode initialize request (initialize): %s\n",
				err,
			)
			return fmt.Errorf(
				"failed to decode initialize request (initialize): %w",
				err,
			)
		}
		s.Logger.Println("received initialize request: ", request)
		s.Logger.Printf(
			"Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version,
		)
		msg := lsp.NewInitializeResponse(request.ID)
		if err = s.writeResponse(msg); err != nil {
			s.Logger.Print(
				fmt.Errorf(
					"failed to write a response: %v the response that failed to write: %s",
					err,
					contents,
				),
			)
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "initialized":
		var request lsp.InitializedParamsRequest
		if err = json.Unmarshal([]byte(contents), &request); err != nil {
			s.Logger.Printf(
				"failed to decode initialized request (initialized): %s\n",
				err,
			)
			return fmt.Errorf(
				"failed to decode initialized request (initialized): %w",
				err,
			)
		}
		s.Logger.Println("received initialized request: ", request)
		msg := lsp.NewInitializedParamsResponse(*request.ID)
		if err = s.writeResponse(msg); err != nil {
			s.Logger.Print(
				fmt.Errorf(
					"failed to write a response: %v the response that failed to write: %s",
					err,
					contents,
				),
			)
			return fmt.Errorf("failed to write response: %w", err)
		}
		enc, err := rpc.EncodeMessage(msg)
		if err != nil {
			s.Logger.Printf("failed to encode message: %s\n", err)
			return fmt.Errorf("failed to encode message: %w", err)
		}
		s.Logger.Println("Received message (initialized) and replied: ", enc)
	case "textDocument/didClose":
		var request lsp.DidCloseTextDocumentParamsNotification
		if err = json.Unmarshal([]byte(contents), &request); err != nil {
			s.Logger.Printf(
				"failed to decode didClose request (didClose): %s\n",
				err,
			)
			return fmt.Errorf(
				"failed to decode didClose request (didClose): %w",
				err,
			)
		}
		s.Logger.Println("received didClose request: ", request)
		msg := lsp.NewDidCloseTextDocumentParamsNotification()
		if err = s.writeResponse(msg); err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err = json.Unmarshal(contents, &request); err != nil {
			return fmt.Errorf(
				"failed to decode didOpen request (textDocument/didOpen): %w",
				err,
			)
		}
		s.Logger.Println("received didOpen request: ", request)
		diagnostics := s.State.OpenDocument(
			request.Params.TextDocument.URI,
			request.Params.TextDocument.Text,
		)
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
		if err = s.writeResponse(alert); err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err = json.Unmarshal(contents, &request); err != nil {
			s.Logger.Printf("failed to decode didChange request: %s\n", err)
			return fmt.Errorf("failed to decode didChange request: %w", err)
		}
		for _, change := range request.Params.ContentChanges {
			diagnostics := s.State.UpdateDocument(
				request.Params.TextDocument.URI,
				change.Text,
			)
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
			if err = s.writeResponse(alert); err != nil {
				s.Logger.Printf(
					"failed to write a response: %s\nthe response that failed to write: %s\n",
					err,
					contents,
				)
				return fmt.Errorf("failed to write response: %w", err)
			}
			s.Logger.Println(
				"Received message (textDocument/didChange) and replied: ",
				alert,
			)
		}
	case "textDocument/hover":
		s.Logger.Println("Received hover message: ", contents)
		var request lsp.HoverRequest
		if err = json.Unmarshal(contents, &request); err != nil {
			s.Logger.Printf("textDocument/hover: %s", err)
			return fmt.Errorf(
				"failed unmarshal of hover request (textDocument/hover): %w",
				err,
			)
		}
		response := s.State.Hover(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		)
		if err = s.writeResponse(response); err != nil {
			s.Logger.Printf(
				"failed to write a response: %s\nthe response that failed to write: %s\n",
				err,
				contents,
			)
			return fmt.Errorf("failed to write response: %w", err)
		}
		s.Logger.Println(
			"Received message (textDocument/didChange) and replied: ",
			contents,
		)
	case "textDocument/definition":
		s.Logger.Println("Received definition message: ", contents)
		var request lsp.DefinitionRequest
		if err = json.Unmarshal(contents, &request); err != nil {
			return s.logE("textDocument/definition", string(contents),
				fmt.Errorf(
					"failed unmarshal of definition request (textDocument/definition): %w",
					err,
				))
		}
		response := s.State.Definition(
			request.ID,
			request.Params.TextDocument.URI,
			request.Params.Position,
		)
		if err = s.writeResponse(response); err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
		s.Logger.Println(
			"Received message (textDocument/definition) and replied: ",
			response,
		)
	case "textDocument/codeAction":
		s.Logger.Println("Received codeAction message: ", contents)
		var request lsp.CodeActionRequest
		if err = json.Unmarshal(contents, &request); err != nil {
			s.Logger.Printf("textDocument/codeAction: %s", err)
			return fmt.Errorf(
				"failed unmarshal of codeAction request (textDocument/codeAction): %w",
				err,
			)
		}
		response := s.State.TextDocumentCodeAction(
			request.ID,
			request.Params.TextDocument.URI,
		)
		if err = s.writeResponse(response); err != nil {
			s.Logger.Printf(
				"failed to write a response: %s\nthe response that failed to write: %s\n",
				err,
				contents,
			)
			return fmt.Errorf("failed to write response: %w", err)
		}
		s.Logger.Println(
			"Received message (textDocument/codeAction) and replied: ",
			response,
		)
	case "textDocument/completion":
		s.Logger.Println("Received completion message: ", contents)
		var request lsp.CompletionRequest
		if err = json.Unmarshal(contents, &request); err != nil {
			s.Logger.Printf("textDocument/codeAction: %s", err)
			return fmt.Errorf(
				"failed unmarshal of completion request (textDocument/completion): %w",
				err,
			)
		}
		response := s.State.TextDocumentCompletion(
			request.ID,
			&request.Params.TextDocument,
			&request.Params.Position,
		)
		if err = s.writeResponse(response); err != nil {
			s.Logger.Printf(
				"failed to write a response: %sthe response that failed to write: %s",
				err,
				contents,
			)
			return fmt.Errorf("failed to write response: %w", err)
		}
		s.Logger.Println(
			"Received message (textDocument/completion) and replied: ",
			response,
		)
	default:
		s.Logger.Println("Unknown method: ", method)
		return nil
	}

	return nil
}

func (s *Root) logF(method string, msg string) {
	s.Logger.Printf("received message (%s): %s\n", method, msg)
}

// logE logs an error for a given method and message
func (s *Root) logE(method string, msg string, err error) error {
	s.Logger.Printf("failed to handle message (%s): %s\n", method, msg)
	return fmt.Errorf("failed to handle message (%s): %w", method, err)
}
