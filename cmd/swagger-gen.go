package cmd

import (
	"fmt"
	"github.com/agorago/wego/fw"
	"log"
	"os"
	"text/template"
)

func swaggergen(wego fw.RegistrationService,service string, templateFile string, targetFile string) error {
	sd, err := wego.FindServiceDescriptor(service)
	if err != nil {
		return err
	}
	tpl, err := template.ParseFiles(templateFile)
	if err != nil {
		fmt.Printf("uh oh problem with template.err = %s\n", err.Error())
		return err
	}

	f, err := os.Create(targetFile)
	if err != nil {
		fmt.Printf("Cannot open %s for writing. Error = %s\n", targetFile, err.Error())
		return err
	}
	err = tpl.Execute(f, sd)
	if err != nil {
		fmt.Printf("Error in writing the template to file %s. Error = %s\n", targetFile, err.Error())
		return err
	}
	return nil
}

// main - this main will need to be invoked by a service after it first loaded its WeGO configurations
// this builds the swagger docs for a specified service that was configured in WeGO
func SwaggerMain(wego fw.RegistrationService){
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s service-name template-file target-file", os.Args[0])
		os.Exit(1)
	}
	serviceName := os.Args[1]
	templateFile := os.Args[2]
	targetFile := os.Args[3]
	err := swaggergen(wego,serviceName, templateFile, targetFile)
	if err != nil {
		log.Fatalf("Cannot generate the file. Error = %s\n", err)
		os.Exit(2)
	}
}
