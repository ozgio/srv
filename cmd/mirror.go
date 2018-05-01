package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/ozgio/srv/pkg/server"
	"github.com/spf13/cobra"
)

const defaultJsonEnabled = false

type mirrorsCmdFlags struct {
	json bool
}

func NewMirrorCommand() *cobra.Command {
	var flags mirrorsCmdFlags

	cmd := &cobra.Command{
		Use:   "mirror",
		Short: "Prints incoming requests",
		Long: `|
Send request as response in text or json format

Examples:
	srv mirror
	srv mirror --port=80 --json
`,
		Run: func(cmd *cobra.Command, args []string) {
			runMirrorCmd(flags, args)
		},
	}

	cmd.Flags().BoolVarP(&flags.json, "json", "j", defaultJsonEnabled, "Response format")

	return cmd
}

func runMirrorCmd(flags mirrorsCmdFlags, args []string) {

	http.HandleFunc("/", mirrorHandlerFunc(flags.json))
	err := server.ListenAndServe(host, port, certFile, keyFile, nil)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
}

func mirrorHandlerFunc(json bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.RequestURI())
		var dump []byte
		var err error
		if json {
			dump, err = createJSON(r)
		} else {
			dump, err = createText(r)
		}
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s", dump)
	}
}

func createText(req *http.Request) ([]byte, error) {
	return httputil.DumpRequest(req, true)
}

func createJSON(req *http.Request) ([]byte, error) {

	var err error
	data := map[string]interface{}{}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	data["body"] = string(body)
	data["uri"] = req.RequestURI
	data["method"] = req.Method
	data["protocol"] = req.Proto
	headers := map[string]string{}
	for name, value := range req.Header {
		headers[name] = strings.Join(value, ", ")
	}
	if req.Close {
		headers["Connection"] = "close"
	}
	if len(req.TransferEncoding) > 0 {
		headers["Transfer-Encoding"] = strings.Join(req.TransferEncoding, ", ")
	}
	data["headers"] = headers

	return json.Marshal(data)

}
