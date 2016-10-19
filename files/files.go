package files

import "os/user"

// InHomeDir returns the path to a file inside the home directory.
func InHomeDir(filename string) (string, error) {
	usr, err := user.Current()
	return inDirOfUser(filename, usr, err)
}

func inDirOfUser(filename string, usr *user.User, err error) (string, error) {
	if err != nil {
		return "", err
	}
	return usr.HomeDir + "/" + filename, nil
}
