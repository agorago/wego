package cmd

import (
	"log"
	"net/http"

	bplusHTTP "gitlab.intelligentb.com/devops/bplus/http"
)

// Serve - start a server for serving HTTP requests
func Serve() {
	log.Fatal(http.ListenAndServe(":5000", bplusHTTP.HTTPHandler))
}
