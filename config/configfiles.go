package config

import (
	"os/user"
	"path"

	"github.com/bpicode/fritzctl/internal/errors"
)

var (
	// Version defines the version of the application.
	Version = "unknown"
	// Revision is the hash in VCS (git commit).
	Revision = "unknown"
)

func homeDirOf(userSupplier func() (*user.User, error)) func(filename string) (string, error) {
	return func(filename string) (string, error) {
		usr, err := userSupplier()
		if err != nil {
			return "", errors.Wrapf(err, "cannot determine current user")
		}
		return path.Join(usr.HomeDir, filename), nil
	}
}
