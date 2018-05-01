package cmd

import (
	"fmt"
	"os"

	"github.com/ozgio/srv/pkg/tls"
	"github.com/spf13/cobra"
)

type generateCmdFlags struct {
}

func NewGenerateCommand() *cobra.Command {
	var flags generateCmdFlags

	var cmd = &cobra.Command{
		Use:   "generate",
		Short: "Generates key and certs files for https server",
		Long: `|
Generates key and certs files for https server. Keep in mind that these are 
only meant for development.

Examples:
	srv generate --key=key.pem --cert=cert.pem
`,
		Run: func(cmd *cobra.Command, args []string) {
			runGenerateCmd(flags, args)
		},
	}

	return cmd
}

func runGenerateCmd(flags generateCmdFlags, args []string) {
	keyFile := keyFile
	certFile := certFile

	if keyFile == "" {
		keyFile = "key.pem"
	}
	if certFile == "" {
		certFile = "cert.pem"
	}

	for _, file := range []string{keyFile, certFile} {
		exists, err := isFileExist(file)
		if err != nil {
			fmt.Printf("error on reading %s: %s\n", file, err.Error())
			return
		}
		if exists {
			fmt.Printf("%s exists. you must delete old one before generating new\n", file)
			return
		}
	}

	opts := tls.Options{
		Host: host,
	}

	fmt.Printf("Generating files %s, %s with these options: %+v\n", keyFile, certFile, opts)

	err := tls.GenerateCertToFiles(keyFile, certFile, opts)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("Files are written succesfully. Now you can use srv commands with these files:\n")
	fmt.Printf("srv files --key=\"%s\" --cert=\"%s\"\n", keyFile, certFile)
}

func isFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
