package server

import (
	"fmt"
	"net/http"
)

func ListenAndServe(port int, certFile, keyFile string, handler http.Handler) error {
	addr := fmt.Sprintf(":%d", port)
	if certFile != "" && keyFile != "" {
		return http.ListenAndServeTLS(addr, "cert.pem", "key.pem", handler)
	}

	return http.ListenAndServe(addr, handler)

}
