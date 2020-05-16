package cmd

import (
	"fmt"
	wego "github.com/agorago/wego"
	"github.com/agorago/wego/config"
	"github.com/agorago/wego/fw"
	"log"
	"net/http"
)

// Serve - start a server for serving HTTP requests
func Serve(initializers ...fw.Initializer) {
	_,httphandler,err := InitApp(initializers...)
	if err != nil {
		panic("Cannot start server. Error = " + err.Error())
	}
	ServeHandle(httphandler)
}

// Serve - start a server for serving HTTP requests
func ServeHandle(httphandler http.Handler ) {
	a := fmt.Sprintf(":%s", config.Value("wego.port"))
	log.Printf("Starting server at address %s", a)
	log.Fatal(http.ListenAndServe(a, httphandler))
}

func InitApp(initializers ...fw.Initializer)(fw.CommandCatalog,http.Handler,error){
	commandCatalog,err := fw.MakeInitializedCommandCatalog(initializers...)
	if err != nil {
		return nil,nil,err
	}

	httphandler,err := wego.GetHTTPHandler(commandCatalog)
	if err != nil {
		return commandCatalog,nil,err
	}
	return commandCatalog,httphandler,nil
}