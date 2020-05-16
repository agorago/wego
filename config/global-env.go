package config

import (
	"os"
)

// Constants for the variables used in this file
const (
	CONFIGPATHVAR        = "CONFIGPATH"
	DEFAULTLANGUAGEVAR   = "wego.default_language"
	ENVVAR               = "ENV"
	ETCD_ENDPOINTVAR     = "wego.etcd_endpoint"
	ETCD_POLLINGDELAYVAR = "wego.etcd_polling_delay"
	APPLICATIONNAMEVAR   = "wego.application_name"
)

func GetEtcdEndPoint() string {
	return Value(ETCD_ENDPOINTVAR)
}

func GetEtcdPollingDelay() int {
	return IntValue(ETCD_ENDPOINTVAR)
}

// GetConfigPath - return the path where all config files are kept
func GetConfigPath() string {
	return os.Getenv(CONFIGPATHVAR)
}

// GetDefaultLanguage - returns the default language to be used if language is
// not specified by the end user
func GetDefaultLanguage() string {
	lang := Value(DEFAULTLANGUAGEVAR)
	if lang == "" {
		return "en-US"
	}
	return lang
}

func GetEnv() string {
	env := os.Getenv(ENVVAR)
	if env == "" {
		return "dev"
	}
	return env
}

func GetApplicationName() string {
	app := Value(APPLICATIONNAMEVAR)
	if app == "" {
		return "wego"
	}
	return app
}
