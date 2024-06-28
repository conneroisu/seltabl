package cmds

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/rpc"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// LSPHandler is a struct for the LSP server
type LSPHandler func(ctx context.Context, writer *io.Writer, state *analysis.State, msg []byte) error

// NewLSPCmd creates a new command for the lsp subcommand
func NewLSPCmd(ctx context.Context, writer io.Writer, handle LSPHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lsp", // the name of the command
		Short: "A command line tooling for package that parsing html tables and elements into structs",
		Long: `
CLI and Language Server for the seltabl package.

Language server provides completions, hovers, and code actions for seltabl defined structs.
	
CLI provides a command line tool for verifying, linting, and reporting on seltabl defined structs.
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SetOut(writer)
			cmd.SetContext(ctx)
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Split(rpc.Split)
			state, err := analysis.NewState()
			if err != nil {
				return fmt.Errorf("failed to create state: %w", err)
			}
			eg := errgroup.Group{}
			for scanner.Scan() {
				eg.Go(func() error {
					msgCtx, cancel := context.WithCancel(ctx)
					defer cancel()
					msg := scanner.Bytes()
					message, err := rpc.DecodeMessage(msg)
					if err != nil {
						return fmt.Errorf("failed to decode message: %w", err)
					}
					if message.ID != nil {
						state.Ctxs.Put(*message.ID, analysis.ContextUnit{Ctx: msgCtx, Cancel: cancel})
						if message.Method == "$/cancelRequest" {
							unit, exists := state.Ctxs.Get(*message.ID)
							if !exists {
								return fmt.Errorf("no context found for id: %d", *message.ID)
							}
							unit.Cancel()
							state.Ctxs.Remove(*message.ID)
						}
					}
					content := bytes.ReplaceAll(message.Content, []byte("\r\n"), []byte(""))
					content = bytes.ReplaceAll(content, []byte("\\"), []byte(""))
					content = bytes.ReplaceAll(content, []byte("/"), []byte(""))
					state.Logger.Printf("received message (%s): %s", message.Method, string(content))
					err = handle(ctx, &writer, &state, msg)
					if err != nil {
						return fmt.Errorf("failed to handle message: %w", err)
					}
					return nil
				})
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("scanner error: %w", err)
				}
				if err := eg.Wait(); err != nil {
					return fmt.Errorf("error group error: %w", err)
				}
			}
			return nil
		},
	}
	return cmd
}

// HandleMessage handles a message sent from the client to the language server.
// It parses the message and returns with a response.
func HandleMessage(
	ctx context.Context,
	writer *io.Writer,
	state *analysis.State,
	msg []byte,
) error {
	message, err := rpc.DecodeMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to decode message: %w", err)
	}
	mux := lsp.DefaultMux
	addRoutes(ctx, mux, state)
	switch message.Method {
	case "initialize":
		var request lsp.InitializeRequest
		if err = json.Unmarshal([]byte(message.Content), &request); err != nil {
			return fmt.Errorf("decode initialize request (initialize) failed: %w", err)
		}
		response := lsp.NewInitializeResponse(request.ID)
		err = WriteResponse(ctx, state, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write (initialize) response: %w", err)
		}
	case "initialized":
		var request lsp.InitializedParamsRequest
		if err = json.Unmarshal([]byte(message.Content), &request); err != nil {
			return fmt.Errorf("decode (initialized) request failed: %w", err)
		}
	case "textDocument/didOpen":
		var request lsp.NotificationDidOpenTextDocument
		if err = json.Unmarshal(message.Content, &request); err != nil {
			return fmt.Errorf(
				"decode (textDocument/didOpen) request failed: %w",
				err,
			)
		}
		diagnostics, err := state.OpenDocument(
			ctx,
			request.Params.TextDocument.URI,
			&request.Params.TextDocument.Text,
		)
		if err != nil {
			return fmt.Errorf("failed to open document: %w", err)
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
		err = WriteResponse(ctx, state, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "textDocument/didClose":
		var request lsp.DidCloseTextDocumentParamsNotification
		if err = json.Unmarshal([]byte(message.Content), &request); err != nil {
			return fmt.Errorf("decode (didClose) request failed: %w", err)
		}
		state.Documents[request.Params.TextDocument.URI] = ""
	case "textDocument/completion":
		var request lsp.CompletionRequest
		err = json.Unmarshal(message.Content, &request)
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
		err = WriteResponse(ctx, state, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		err = json.Unmarshal(message.Content, &request)
		if err != nil {
			return fmt.Errorf("decode (textDocument/didChange) request failed: %w", err)
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
		err = WriteResponse(ctx, state, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "textDocument/hover":
		var request lsp.HoverRequest
		err = json.Unmarshal(message.Content, &request)
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
		err = WriteResponse(ctx, state, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		err = json.Unmarshal(message.Content, &request)
		if err != nil {
			return fmt.Errorf("failed to unmarshal of codeAction request (textDocument/codeAction): %w", err)
		}
		response, err := state.TextDocumentCodeAction(
			request.ID,
			request.Params.TextDocument.URI,
		)
		if err != nil {
			return fmt.Errorf("failed to get code actions: %w", err)
		}
		err = WriteResponse(ctx, state, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "textDocument/didSave":
	case "shutdown":
		var request lsp.ShutdownRequest
		err = json.Unmarshal([]byte(message.Content), &request)
		if err != nil {
			return fmt.Errorf("decode (shutdown) request failed: %w", err)
		}
		response := lsp.NewShutdownResponse(request, nil)
		err = WriteResponse(ctx, state, writer, response)
		if err != nil {
			return fmt.Errorf("write (shutdown) response failed: %w", err)
		}
		os.Exit(0)
	case "$/cancelRequest":
		var request lsp.CancelRequest
		err = json.Unmarshal(message.Content, &request)
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
		err = WriteResponse(ctx, state, writer, response)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case "exit":
		os.Exit(0)
		return nil
	default:
		state.Logger.Printf("unknown method: %s", message.Method)
	}
	return nil
}
