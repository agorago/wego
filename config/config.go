package config

import (
	"os"
)

// Constants for the variables used in this file
const (
	CONFIGPATHVAR      = "CONFIGPATH"
	DEFAULTLANGUAGEVAR = "DEFAULTLANGUAGE"
)

// GetConfigPath - return the path where all config files are kept
func GetConfigPath() string {
	return os.Getenv(CONFIGPATHVAR)
}

// GetDefaultLanguage - returns the default language to be used if language is
// not specified by the end user
func GetDefaultLanguage() string {
	lang := os.Getenv(DEFAULTLANGUAGEVAR)
	if lang == "" {
		return "en-US"
	}
	return lang
}
