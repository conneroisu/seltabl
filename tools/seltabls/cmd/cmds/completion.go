package cmds

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// NewCompletionCmd returns the completion command
func NewCompletionCmd(
	ctx context.Context,
	w io.Writer,
	r io.Reader,
) *cobra.Command {
	var completionCmd = &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generates completion scripts for docuvet-tools for your shell",
		Long: `To load completions:
	Bash:

	  $ source <(%[1]s completion bash)

	  # To load completions for each session, execute once:
	  # Linux:
	  $ %[1]s completion bash > /etc/bash_completion.d/%[1]s
	  # macOS:
	  $ %[1]s completion bash > $(brew --prefix)/etc/bash_completion.d/%[1]s

	Zsh:

	  # If shell completion is not already enabled in your environment,
	  # you will need to enable it.  You can execute the following once:

	  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

	  # To load completions for each session, execute once:
	  $ %[1]s completion zsh > "${fpath[1]}/_%[1]s"

	  # You will need to start a new shell for this setup to take effect.

	fish:

	  $ %[1]s completion fish | source

	  # To load completions for each session, execute once:
	  $ %[1]s completion fish > ~/.config/fish/completions/%[1]s.fish

	PowerShell:

	  PS> %[1]s completion powershell | Out-String | Invoke-Expression

	  # To load completions for every new session, run:
	  PS> %[1]s completion powershell > %[1]s.ps1
	  # and source this file from your PowerShell profile.
	`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cmd.SetIn(r)
			cmd.SetOut(w)
			_, cancel := context.WithCancel(ctx)
			defer cancel()
			if len(args) == 0 {
				err = cmd.Help()
				if err != nil {
					return fmt.Errorf("failed to get help: %w", err)
				}
				return
			}
			switch args[0] {
			case "bash":
				err = cmd.Root().GenBashCompletion(w)
			case "zsh":
				err = cmd.Root().GenZshCompletion(w)
			case "fish":
				err = cmd.Root().GenFishCompletion(w, true)
			case "powershell":
				err = cmd.Root().
					GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			}
			if err != nil {
				return fmt.Errorf("failed to generate completions: %w", err)
			}
			return nil
		},
	}
	return completionCmd
}
