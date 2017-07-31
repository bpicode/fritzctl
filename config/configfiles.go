package config

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/bpicode/fritzctl/functional"
)

var (
	// ApplicationName denotes the name of the application.
	ApplicationName = "fritzctl"
	// Version defines the version of the application.
	Version = "unknown"
	// ConfigFilename defines the filename of the configuration file.
	ConfigFilename = "fritzctl.json"
	// ConfigFilenameHidden defines the filename of the configuration file (hidden).
	ConfigFilenameHidden = "." + ConfigFilename
	// ConfigDir defines the directory of the configuration file.
	ConfigDir = ""
	// DefaultConfigDir is the default directory where the config file resides.
	DefaultConfigDir = "/etc/fritzctl"
)

// FindConfigFile returns the path to the config file.
func FindConfigFile() (string, error) {
	return functional.FirstWithoutError(
		functional.Curry(fmt.Sprintf("%s/%s", ConfigDir, ConfigFilename), accessible),
		functional.Compose(ConfigFilenameHidden, homeDirOf(user.Current), accessible),
		functional.Curry(DefaultConfigFileAbsolute(), accessible),
	)
}

// DefaultConfigFileAbsolute returns the absolute path of the default configuration file.
func DefaultConfigFileAbsolute() string {
	return fmt.Sprintf("%s/%s", DefaultConfigDir, ConfigFilename)
}

func accessible(file string) (string, error) {
	_, err := os.Stat(file)
	return file, err
}

func homeDirOf(userSupplier func() (*user.User, error)) func(filename string) (string, error) {
	return func(filename string) (string, error) {
		usr, err := userSupplier()
		if err != nil {
			return "", fmt.Errorf("cannot determine current user: %s", err)
		}
		return path.Join(usr.HomeDir, filename), nil
	}
}
