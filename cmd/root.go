package cmd

import (
	"github.com/spf13/cobra"
)

const defaultServerPort = 8010
const defaultProtocol = "http"
const defaultHost = "127.0.0.1"

var port int
var certFile string
var keyFile string
var host string

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "srv",
		Short: "srv is a http server for helping development",
		Long:  `http server for helping development`,
	}
	return cmd
}

func Execute() (*cobra.Command, error) {
	rootCmd := NewRootCommand()
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", defaultServerPort, "Port to listen")
	rootCmd.PersistentFlags().StringVarP(&certFile, "host", "l", defaultHost, "Host name or address")
	rootCmd.PersistentFlags().StringVarP(&certFile, "cert", "c", "", "Path to cert file for https server")
	rootCmd.PersistentFlags().StringVarP(&keyFile, "key", "k", "", "Path to key file for https server")

	rootCmd.AddCommand(NewFilesCommand())
	rootCmd.AddCommand(NewMirrorCommand())

	err := rootCmd.Execute()

	return rootCmd, err
}
