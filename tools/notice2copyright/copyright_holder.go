package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var copyrightRegex = regexp.MustCompile(`^>*\s*[Cc]opyright\s(\([Cc]\)|©)*(?P<CLINE>.*)`)

var copyrightNoAuthorRegex = regexp.MustCompile(`([Cc]opyright notice|[Cc]opyright license to reproduce)`)

func findCopyrightHolders(projects []project, dir string) []project {
	withCpHolders := make([]project, len(projects))
	for i, p := range projects {
		withCpHolders[i] = findCopyrightHoldersOf(p, dir)
	}
	return withCpHolders
}

func findCopyrightHoldersOf(p project, dir string) project {
	srcDir := filepath.Join(dir, p.dir())
	lf, err := openFirst(srcDir, "LICENSE", "LICENSE.md", "LICENSE.txt", "LICENSE.rst")
	assertOrFatal(err == nil, "unable to open license file for project '%s' (%s): %v", p.name, p.url, err)
	defer lf.Close()
	p.copyrightHolders = readCopyrightHolders(lf)
	return p
}

func openFirst(dir string, names ...string) (*os.File, error) {
	for _, name := range names {
		f, err := os.Open(filepath.Join(dir, name))
		if err == nil {
			return f, nil
		}
	}
	return nil, fmt.Errorf("no license file [%s] found in directory '%s'", names, dir)
}

func readCopyrightHolders(r io.Reader) []string {
	var holders []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		holders = addCopyRightHolder(line, holders)
	}
	if len(holders) == 0 {
		holders = append(holders, "unknown")
	}
	return holders
}

func addCopyRightHolder(text string, holders []string) []string {
	if copyrightRegex.MatchString(text) && !copyrightNoAuthorRegex.MatchString(text) {
		holders = append(holders, parseCopyrightHolder(text))
	}
	return holders
}

func parseCopyrightHolder(text string) string {
	matches := copyrightRegex.FindStringSubmatch(text)
	return strings.TrimSpace(matches[len(matches)-1])
}
