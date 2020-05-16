package nr

import (
	"log"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/agorago/wego/config"
)

var NRApp *newrelic.Application

func init() {
	if !config.BoolValue("wego.new_relic_enabled") {
		return
	}
	var err error
	NRApp, err = newrelic.NewApplication(
		newrelic.ConfigAppName(config.GetApplicationName()),
		newrelic.ConfigLicense(config.Value("wego.new_relic_license_key")),
		newrelic.ConfigInfoLogger(os.Stdout),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if nil != err {
		log.Printf("Cannot create the new relic agent. Error = %s\n", err.Error())
	}
}
