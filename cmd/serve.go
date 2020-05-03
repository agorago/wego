package cmd

import (
	"fmt"
	"github.com/agorago/wego/config"
	"log"
	"net/http"

	wegohttp "github.com/agorago/wego/http"
)

// Serve - start a server for serving HTTP requests
func Serve() {
	a := fmt.Sprintf(":%s", config.Value("bplus.port"))
	log.Printf("Starting server at address %s", a)
	log.Fatal(http.ListenAndServe(a, wegohttp.HTTPHandler))
}
