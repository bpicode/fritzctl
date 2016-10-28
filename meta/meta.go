package meta

import (
	"fmt"

	"os"

	"github.com/bpicode/fritzctl/files"
	"github.com/bpicode/fritzctl/functional"
)

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
	return functional.FirstWithoutError(
		functional.Compose(ConfigFilename, files.InHomeDir, accessible),
		functional.Curry(fmt.Sprintf("%s/%s", "/etc/fritzctl", ConfigFilename), accessible),
		functional.Curry(fmt.Sprintf("%s/%s", ConfigDir, ConfigFilename), accessible),
	)
}

func accessible(file string) (string, error) {
	_, err := os.Stat(file)
	return file, err
}
