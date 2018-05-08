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
		Short: "srv is a http server for helping development",
		Long:  `http server for helping development`,
	}
	return cmd
}

// Execute parses flags and runs commands. It is the starting point of the
// application. Returns the command which is the root of all sub commands
func Execute() (*cobra.Command, error) {
	rootCmd := NewRootCommand()
	rootCmd.PersistentFlags().IntVarP(&defaultRootCmdFlags.port, "port", "p", defaultServerPort, "Port to listen")
	rootCmd.PersistentFlags().StringVarP(&defaultRootCmdFlags.host, "host", "o", defaultHost, "Host name or address")
	rootCmd.PersistentFlags().StringVarP(&defaultRootCmdFlags.certFile, "cert", "c", "", "Path to cert file for https server")
	rootCmd.PersistentFlags().StringVarP(&defaultRootCmdFlags.keyFile, "key", "k", "", "Path to key file for https server")

	rootCmd.AddCommand(NewFilesCommand())
	rootCmd.AddCommand(NewMirrorCommand())
	rootCmd.AddCommand(NewGenerateCommand())

	err := rootCmd.Execute()

	return rootCmd, err
}
