package cmd

import (
	"fmt"
	"os"

	"github.com/ozgio/srv/pkg/tls"
	"github.com/spf13/cobra"
)

type generateCmdFlags struct {
}

// NewGenerateCommand creates new "generate" subcommand
func NewGenerateCommand() *cobra.Command {
	var flags generateCmdFlags

	var cmd = &cobra.Command{
		Use:   "generate",
		Short: "Generates key and certs files for https server",
		Long: `
Generates key and cert files for https server. Keep in mind that these are 
only meant for development. These files can be used with "srv files" or 
"srv mirror" commands later

Examples:
	srv generate --key=key.pem --cert=cert.pem

	#create https server using generated certificate
	srv files --cert=cert.pem --key=key.pem
`,
		Run: func(cmd *cobra.Command, args []string) {
			runGenerateCmd(flags, args)
		},
	}

	return cmd
}

func runGenerateCmd(flags generateCmdFlags, args []string) {
	keyFile := defaultRootCmdFlags.keyFile
	certFile := defaultRootCmdFlags.certFile

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
		Host: defaultRootCmdFlags.host,
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
