package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	assertOrFatal(len(os.Args) >= 3, "usage notice2copyright <root project directory> <root project license>")
	root := parseRoot(os.Args[1])
	deps := parseNotice(filepath.Join(os.Args[1], "NOTICE"))
	deps = findCopyrightHolders(deps, filepath.Dir(os.Args[1]))
	writer := copyrightWriter{root: root, deps: deps}
	writer.writeTo(os.Stdout)
}

func parseRoot(dir string) project {
	abs, err := filepath.Abs(dir)
	assertOrFatal(err == nil, "cannot determine rot project dir: %v", err)
	dirs := strings.Split(abs, string(filepath.Separator))
	root := parseRootFromDirHierarchy(dirs)
	return findCopyrightHoldersOf(root, dir)
}

func parseRootFromDirHierarchy(dirs []string) project {
	assertOrFatal(len(dirs) >= 3, "cannot determine root project: project should reside in a folder like '.../github.com/user/project', got instead %s", dirs)
	name, user, host := dirs[len(dirs)-1], dirs[len(dirs)-2], dirs[len(dirs)-3]
	url := fmt.Sprintf("https://%s/%s/%s", host, user, name)
	p := project{
		url:     url,
		name:    name,
		license: licenses[os.Args[2]],
		isRoot:  true,
	}
	return p
}
