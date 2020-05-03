package cmd

import (
	"fmt"
	"github.com/agorago/wego/fw"
	"log"
	"os"
	"text/template"
)

func swaggergen(service string, templateFile string, targetFile string) error {
	sd, err := fw.FindServiceDescriptor(service)
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

// main - this main will need to be invoked by a service after it first loaded its BPLUS configurations
// this builds the swagger docs for a specified service that was configured in BPlus
func SwaggerMain() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s service-name template-file target-file", os.Args[0])
		os.Exit(1)
	}
	serviceName := os.Args[1]
	templateFile := os.Args[2]
	targetFile := os.Args[3]
	err := swaggergen(serviceName, templateFile, targetFile)
	if err != nil {
		log.Fatalf("Cannot generate the file. Error = %s\n", err)
		os.Exit(2)
	}
}
