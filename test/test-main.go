package test

import (
	"flag"
	"os"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/MenaEnergyVentures/bplus/cmd"
	_ "github.com/MenaEnergyVentures/bplus/http" // ensure http is registered
)

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

// BDDSuite - the actual bdd suite that contains the step definitions
type BDDSuite func(*godog.Suite)

// BDD - the method that invokes the goDog BDD suite of tests
func BDD(m *testing.M, bddsuite BDDSuite) {
	flag.Parse()
	opt.Paths = flag.Args()

	go cmd.Serve() // this is important. Else the server wont start. It is also important that
	// the server is not running in the foreground since we need to initiate the tests after this

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		bddsuite(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}
