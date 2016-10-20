package meta

import "github.com/bpicode/fritzctl/files"

var (
	// ApplicationName denotes the name of the application.
	ApplicationName = "fritzctl"
	// Version defines the version of the application.
	Version = "1.0.0"
	// ConfigFilename defines the filename of the configuration file.
	ConfigFilename = "fritzctl.json"
	// ConfigDir defines the directory of the configuration file.
	ConfigDir = ""
)

// ConfigFile returns the path to the config file.
func ConfigFile() (string, error) {
	if ConfigDir == "" {
		return files.InHomeDir(ConfigFilename)
	}
	return ConfigDir + "/" + ConfigFilename, nil
}
