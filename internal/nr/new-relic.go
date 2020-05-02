package nr

import (
	"log"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
	"gitlab.intelligentb.com/devops/bplus/config"
)

var NRApp *newrelic.Application

func init() {
	if !config.BoolValue("bplus.new_relic_enabled") {
		return
	}
	var err error
	NRApp, err = newrelic.NewApplication(
		newrelic.ConfigAppName(config.GetApplicationName()),
		newrelic.ConfigLicense(config.Value("bplus.new_relic_license_key")),
		newrelic.ConfigInfoLogger(os.Stdout),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if nil != err {
		log.Printf("Cannot create the new relic agent. Error = %s\n", err.Error())
	}
}
