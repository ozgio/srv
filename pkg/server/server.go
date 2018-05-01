package server

import (
	"fmt"
	"net/http"
)

// ListenAndServe starts http or https server based on parameters certFile and
// keyFile
//
// Returns whether the server is https or not as bool and the error
// http.ListenAndServer* returns
func ListenAndServe(host string, port int, certFile, keyFile string, handler http.Handler) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	var proto = "http"
	if certFile != "" && keyFile != "" {
		proto = "https"
	}

	fmt.Printf("Server started running at %s://%s\n", proto, addr)
	if proto == "https" {
		return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
	}

	return http.ListenAndServe(addr, handler)

}
