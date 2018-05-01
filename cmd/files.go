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

func NewFilesCommand() *cobra.Command {
	var flags filesCmdFlags

	var cmd = &cobra.Command{
		Use:   "files",
		Short: "Static file server",
		Long: `|
Static file server with directory listing

Examples:
	srv files
	srv files --port=80 --root=~/files
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

	fmt.Printf("Static file server running at %s://%s:%d\n", defaultProtocol, "localhost", port)
	err = server.ListenAndServe(port, certFile, keyFile, http.FileServer(root))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

}
