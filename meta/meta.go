package meta

import (
	"errors"

	"fmt"

	"os"

	"github.com/bpicode/fritzctl/files"
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
	return firstWithoutError(
		compose(ConfigFilename, files.InHomeDir, accessible),
		curry(fmt.Sprintf("%s/%s", "/etc/fritzctl", ConfigFilename)+ConfigFilename, accessible),
		curry(fmt.Sprintf("%s/%s", ConfigDir, ConfigFilename), accessible),
	)
}

func firstWithoutError(fcs ...func() (string, error)) (string, error) {
	var ret string
	var err error
	var errs []error
	for _, f := range fcs {
		ret, err = f()
		if err == nil {
			return ret, nil
		}
		errs = append(errs, err)
	}
	return "", errors.New(fmt.Sprint(errs))
}

func accessible(file string) (string, error) {
	_, err := os.Stat(file)
	return file, err
}

func curry(arg string, f func(string) (string, error)) func() (string, error) {
	return func() (string, error) {
		return f(arg)
	}
}

func compose(arg0 string, fcs ...func(string) (string, error)) func() (string, error) {
	return func() (string, error) {
		arg := arg0
		for _, f := range fcs {
			var err error
			arg, err = f(arg)
			if err != nil {
				return "", err
			}
		}
		return arg, nil
	}
}
