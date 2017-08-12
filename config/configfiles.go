package config

import (
	"fmt"
	"os"
	"os/user"
	"path"
)

var (
	// Version defines the version of the application.
	Version = "unknown"
	// Filename defines the filename of the configuration file.
	Filename = "fritzctl.json"
	// filenameHidden defines the filename of the configuration file (hidden).
	filenameHidden = "." + Filename
	// Dir defines the directory of the configuration file.
	Dir = ""
	// DefaultDir is the default directory where the config file resides.
	DefaultDir = "/etc/fritzctl"
)

// FindConfigFile returns the path to the config file.
func FindConfigFile() (string, error) {
	return firstWithoutError(
		curry(fmt.Sprintf("%s/%s", Dir, Filename), accessible),
		compose(filenameHidden, homeDirOf(user.Current), accessible),
		curry(DefaultConfigFileAbsolute(), accessible),
	)
}

// DefaultConfigFileAbsolute returns the absolute path of the default configuration file.
func DefaultConfigFileAbsolute() string {
	return fmt.Sprintf("%s/%s", DefaultDir, Filename)
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
