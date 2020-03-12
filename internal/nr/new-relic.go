package nr

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"gitlab.intelligentb.com/devops/bplus/config"
	"log"
	"os"
)

var NRApp *newrelic.Application
func init() {
	var err error
	NRApp, err = newrelic.NewApplication(
		newrelic.ConfigAppName(config.GetApplicationName()),
		newrelic.ConfigLicense(config.Value("new_relic_license_key")),
		newrelic.ConfigDebugLogger(os.Stdout),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if nil != err {
		log.Printf("Cannot create the new relic agent. Error = %s\n", err.Error())
	}
}
