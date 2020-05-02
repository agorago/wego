package cmd

import (
	"fmt"
	"gitlab.intelligentb.com/devops/bplus/config"
	"log"
	"net/http"

	bplusHTTP "gitlab.intelligentb.com/devops/bplus/http"
)

// Serve - start a server for serving HTTP requests
func Serve() {
	a := fmt.Sprintf(":%s", config.Value("bplus.port"))
	log.Printf("Starting server at address %s", a)
	log.Fatal(http.ListenAndServe(a, bplusHTTP.HTTPHandler))
}
