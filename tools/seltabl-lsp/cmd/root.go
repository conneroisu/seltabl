package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/rpc"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
)

// ReturnCmd returns the command for the root
func (s *Root) ReturnCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "seltabl-lsp", // the name of the command
		Short: "A command line tool for parsing html tables into structs",
		Long: `Language Server for the seltabl package.

Provides completions, hovers, and code actions for seltabl defined structs.
`,
		RunE: func(_ *cobra.Command, _ []string) (err error) {
			s.Logger = getLogger("./seltabl.log")
			s.State.Logger = getLogger("./state.log")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Split(rpc.Split)
			analysis.NewState(os.Getenv)
			go func() {
				if r := recover(); r != nil {
					s.Logger.Printf("Recovered: %v\n", r)
					s.State.Logger.Printf("Recovered: %v\n", r)
				}
			}()
			for scanner.Scan() {
				msg := scanner.Bytes()
				err = s.HandleMessage(msg)
				if err != nil {
					s.Logger.Printf("failed to handle message: %s\n", err)
					s.State.Logger.Printf("failed to handle message: %s\n", err)
					continue
				}
			}
			return nil
		},
	}
}

// Execute runs the root command
func Execute(srv lsp.Server) error {
	cmd := srv.ReturnCmd()
	err := cmd.Execute()
	if err != nil {
		return fmt.Errorf("failed to execute root command: %w", err)
	}
	return nil
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

// Root is the server for the root command
type Root struct {
	lsp.Server
	// State is the State of the server
	State analysis.State
	// Logger is the Logger for the server
	Logger *log.Logger
	// Writer is the Writer for the server
	Writer io.Writer
	// Database is the database for the server
	Database *bun.DB
}

// writeResponse writes a message to the writer
func (s *Root) writeResponse(msg interface{}) error {
	reply, err := rpc.EncodeMessage(msg)
	if err != nil {
		s.Logger.Printf("failed to encode message: %s\n", err)
		return fmt.Errorf("failed to encode message: %w", err)
	}
	res, err := s.Writer.Write([]byte(reply))
	if err != nil {
		s.Logger.Printf("failed to write response: %s\n", err)
		return fmt.Errorf("failed to write response: %w", err)
	}
	if res != len(reply) {
		s.Logger.Printf("failed to write all response: %s\n", err)
		return fmt.Errorf("failed to write all response: %w", err)
	}
	return nil
}
