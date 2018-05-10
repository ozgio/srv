//Package cmd provides cobra commands
package cmd

import (
	"github.com/spf13/cobra"
)

//defaults for flags
const (
	defaultServerPort = 8010
	defaultProtocol   = "http"
	defaultHost       = "127.0.0.1"
)

var version string

// rootCmdFlags represents cli flags for root command
type rootCmdFlags struct {
	port     int
	certFile string
	keyFile  string
	host     string
}

//global varibale for accessing root flags (persistentFlags)
var defaultRootCmdFlags rootCmdFlags

//NewRootCommand creates a new Command as root.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "srv",
		Short: "srv is a HTTP server for helping development",
		Long: `
HTTP server for helping development which can serve static files or prints
incoming requests.

srv supports TLS (https) connection and comes with a self signed TLS cerificate
generator.

Examples:
	
	#serve the files at current working directory on port 8010
	srv files

	#generate self-signed TLS certificates
	srv generate --key=key.pem --cert=cert.pem

	#static file server which supports https
	srv files --port=443 --cert=cert.pem --key=key.pem
`,
	}
	return cmd
}

// Execute parses flags and runs commands. It is the starting point of the
// application. Returns the command which is the root of all sub commands
func Execute(ver string) (*cobra.Command, error) {
	version = ver

	rootCmd := NewRootCommand()
	rootCmd.PersistentFlags().IntVarP(&defaultRootCmdFlags.port, "port", "p", defaultServerPort, "Port to listen")
	rootCmd.PersistentFlags().StringVarP(&defaultRootCmdFlags.host, "host", "o", defaultHost, "Host name or address")
	rootCmd.PersistentFlags().StringVarP(&defaultRootCmdFlags.certFile, "cert", "c", "", "Path to cert file for https server")
	rootCmd.PersistentFlags().StringVarP(&defaultRootCmdFlags.keyFile, "key", "k", "", "Path to key file for https server")

	rootCmd.AddCommand(NewFilesCommand())
	rootCmd.AddCommand(NewMirrorCommand())
	rootCmd.AddCommand(NewGenerateCommand())
	rootCmd.AddCommand(NewVersionCommand(version))

	err := rootCmd.Execute()

	return rootCmd, err
}
