package cmds

import "github.com/spf13/cobra"

// NewVersionCmd returns the version command.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of the seltabls command",
		Long:  `All software has versions. This is seltabls's`,
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Println("seltabls version 0.1.0.5.0.0-beta1.final")
		},
	}
}
