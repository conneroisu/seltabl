package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/internal/config"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/rpc"
	"github.com/spf13/cobra"
)

// ReturnCmd returns the command for the root
func (s *Root) ReturnCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "seltabl", // the name of the command
		Short: "A command line tool for parsing html tables into structs",
		Long: `
CLI and Language Server for the seltabl package.

Language server provides completions, hovers, and code actions for seltabl defined structs.
	
CLI provides a command line tool for verifying, linting, and reporting on seltabl defined structs.
`,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := config.CreateConfigDir()
			if err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}
			s.Config = cfg
			s.State, err = analysis.NewState(s.Config)
			if err != nil {
				return fmt.Errorf("failed to create state: %w", err)
			}
			s.Logger = getLogger(path.Join(s.Config.ConfigPath, "seltabl.log"))
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Split(rpc.Split)
			for scanner.Scan() {
				s.handle(scanner)
			}
			return nil
		},
	}
}

// handle handles a message from the client to the language server.
func (s *Root) handle(scanner *bufio.Scanner) {
	defer func() {
		if err := scanner.Err(); err != nil {
			out := os.Stderr
			fmt.Fprintf(out, "scanner error: %v\n", err)
			s.Logger.Printf("scanner error: %v\n", err)
			s.State.Logger.Printf("scanner error: %v\n", err)
		}
	}()
	msg := scanner.Bytes()
	out := os.Stderr
	err := s.HandleMessage(msg)
	if err != nil {
		fmt.Fprintf(out, "failed to handle message: %s\n", err)
		s.Logger.Printf("failed to handle message: %s\n", err)
		s.State.Logger.Printf("failed to handle message: %s\n", err)
		return
	}
}

// Execute runs the root command
func Execute() error {
	srv := &Root{Writer: os.Stdout}
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
	// Config is the config for the server
	Config *config.Config
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
