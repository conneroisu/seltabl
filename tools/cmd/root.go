package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/conneroisu/seltabl/tools/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/pkg/rpc"
	"github.com/spf13/cobra"
)

// Execute runs the root command
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute root command: %w", err)
	}
	return nil
}

// writeResponse writes a message to the writer
func writeResponse(writer io.Writer, msg any) error {
	reply, err := rpc.EncodeMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}
	res, err := writer.Write([]byte(reply))
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}
	if res != len(reply) {
		return fmt.Errorf("failed to write all response: %w", err)
	}
	return nil
}

// rootCmd is the root command for the command line tool
var rootCmd = &cobra.Command{
	Use:   "seltabl-lsp", // the name of the command
	Short: "A command line tool for parsing html tables into structs",
	Long: `A command line tool for parsing html tables into structs.

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := getLogger("./seltabl.log")
		logger.Println("Starting the server...")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(rpc.Split)
		state := analysis.NewState()

		for scanner.Scan() {
			msg := scanner.Bytes()
			logger.Printf("Received message: %s\n", msg)
			method, content, err := rpc.DecodeMessage(msg)
			if err != nil {
				logger.Printf("failed to decode message: %s\n", err)
				continue
			}
			err = handleMessage(
				logger,
				os.Stdout,
				&state,
				method,
				content,
			)
			if err != nil {
				logger.Printf("failed to handle message: %s\n", err)
				continue
			}
		}
		return nil
	},
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
