package cmd

import (
	"log"
	"net/http"

	bplusHTTP "github.com/MenaEnergyVentures/bplus/http"
)

// Serve - start a server for serving HTTP requests
func Serve() {
	log.Fatal(http.ListenAndServe(":8080", bplusHTTP.HTTPHandler))
}
