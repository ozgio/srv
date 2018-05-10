package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewFilesCommand creates new "files" subcommand
func NewVersionCommand(version string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the version of the program",
		Long: `|
Prints the version of the program

Examples:
	srv version
`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	return cmd
}
