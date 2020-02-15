package config

import (
	"os"
)

// Constants for the variables used in this file
const (
	CONFIGPATHVAR = "CONFIGPATH"
)

// GetConfigPath - return the path where all config files are kept
func GetConfigPath() string {
	return os.Getenv(CONFIGPATHVAR)
}
