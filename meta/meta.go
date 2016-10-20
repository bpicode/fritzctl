package meta

import "github.com/bpicode/fritzctl/files"

var (
	// ApplicationName denotes the name of the application.
	ApplicationName = "fritzctl"
	// Version defines the verison of the application.
	Version = "1.0.0"
	// ConfigFilename defines the filename of the configuration file.
	ConfigFilename = "fritzctl.json"
)

// ConfigFile returns the path to the config file.
func ConfigFile() (string, error) {
	return files.InHomeDir(ConfigFilename)
}
