package cmd

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/ozgio/srv/pkg/server"
	"github.com/spf13/cobra"
)

const defaultRoot = "./"

type filesCmdFlags struct {
	root string
}

// NewFilesCommand creates new "files" subcommand
func NewFilesCommand() *cobra.Command {
	var flags filesCmdFlags

	var cmd = &cobra.Command{
		Use:   "files",
		Short: "Static file server",
		Long: `
Static file server with directory listing

Examples:
	#serve the files at current working directory
	srv files

	#serve the files at ~/files using port 80
	srv files --port=80 --root=~/files

	#create https server using specified pem files
	srv files --port=443 --cert=cert.pem --key=key.pem
`,
		Run: func(cmd *cobra.Command, args []string) {
			runCheckCmd(flags, args)
		},
	}

	cmd.Flags().StringVarP(&flags.root, "root", "r", defaultRoot, "Root path for server")

	return cmd
}

func runCheckCmd(flags filesCmdFlags, args []string) {
	rootPath, err := filepath.Abs(flags.root)
	if err != nil {
		fmt.Printf("Cannot open root path: %s\n", err.Error())
	}

	root := http.Dir(rootPath)

	err = server.ListenAndServe(defaultRootCmdFlags.host, defaultRootCmdFlags.port, defaultRootCmdFlags.certFile, defaultRootCmdFlags.keyFile, http.FileServer(root))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

}
