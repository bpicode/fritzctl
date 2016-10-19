package files

import "os/user"

// InHomeDir returns the path to a file inside the home directory.
func InHomeDir(filename string) (string, error) {
	usr, errUser := user.Current()
	if errUser != nil {
		return "", errUser
	}
	return usr.HomeDir + "/" + filename, nil
}
