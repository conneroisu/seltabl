package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/conneroisu/seltabl/tools/lsp"
	"github.com/conneroisu/seltabl/tools/rpc"
	"github.com/spf13/cobra"
)

// Execute runs the root command
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute root command: %w", err)
	}
	return nil
}

// rootCmd is the root command for the command line tool
var rootCmd = &cobra.Command{
	Use:   "seltabl-lsp", // the name of the command
	Short: "A command line tool for parsing html tables into structs",
	Long: `A command line tool for parsing html tables into structs.

The command line tool is used to parse html tables into structs.

The command line tool can be used to parse html tables into structs.

The command line tool can be used to parse html tables into structs.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := getLogger("./seltabl.log")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(rpc.Split)
		writer := os.Stdout
		for scanner.Scan() {
			msg := scanner.Bytes()
			logger.Printf("Received message: %s\n", msg)
			method, content, err := rpc.DecodeMessage(msg)
			if err != nil {
				logger.Printf("failed to decode message: %s\n", err)
				continue
			}
			handleMessage(logger, writer, method, content)
		}
		return nil
	},
}

// handleMessage handles a message
func handleMessage(
	logger *log.Logger,
	writer io.Writer,
	method string,
	msg []byte,
) {
	logger.Printf("Received message (%s): %s\n", method, msg)
	switch method {
	case "initialize":
		logger.Println("received initialize request")
		var request lsp.InitializeRequest
		if err := json.Unmarshal([]byte(msg), &request); err != nil {
			logger.Printf("failed to decode initialize request: %s\n", err)
			return
		}
		logger.Printf(
			"Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version,
		)
	case "textDocument/didOpen":
		logger.Println("Received didOpen message")
	case "textDocument/didChange":
		logger.Println("Received didChange message")
	case "textDocument/didClose":
		logger.Println("Received didClose message")
	default:
		logger.Printf("Unknown method: %s\n", method)
	}
}

// getLogger returns a logger that writes to a file
func getLogger(fileName string) *log.Logger {
	logFile, err := os.OpenFile(
		fileName,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	return log.New(logFile, "[seltabl-lsp]", log.LstdFlags)
}
