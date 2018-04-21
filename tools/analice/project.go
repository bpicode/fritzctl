package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type project struct {
	root         string
	dependencies []dependency
}

func (p *project) short() string {
	name, _, _ := p.tokens()
	return name
}

func (p *project) url() string {
	name, user, host := p.tokens()
	return fmt.Sprintf("https://%s/%s/%s", host, user, name)
}

func (p *project) tokens() (string, string, string) {
	abs, _ := filepath.Abs(p.root)
	split := strings.Split(abs, string(filepath.Separator))
	padded := append([]string{"", "", ""}, split...)
	return padded[len(padded)-1], padded[len(padded)-2], padded[len(padded)-3]
}

type dependency struct {
	name string
}

func (d *dependency) short() string {
	sp := strings.Split(d.name, "/")
	return sp[len(sp)-1]
}

type projector interface {
	project(dir string) (*project, error)
}

type depLockProjector struct {
}

func (dlp *depLockProjector) project(dir string) (*project, error) {
	join := path.Join(dir, "Gopkg.lock")
	file, err := os.Open(join)
	if err != nil {
		return nil, fmt.Errorf("failed to open Gopkg.lock: %v", err)
	}
	defer file.Close()
	return &project{root: dir, dependencies: dlp.parse(file)}, nil
}

var projectNameRegex = regexp.MustCompile(`^\s*name\s*=\s*"?(?P<NAME>[^"\s]+|\S+)"?\s*`)

func (dlp *depLockProjector) parse(r io.Reader) []dependency {
	s := bufio.NewScanner(r)
	var deps []dependency
	for s.Scan() {
		text := s.Text()
		deps = dlp.condAppendProject(deps, text)
	}
	return deps
}

func (dlp *depLockProjector) condAppendProject(deps []dependency, text string) []dependency {
	if projectNameRegex.MatchString(text) {
		n := dlp.parseName(text)
		deps = append(deps, dependency{name: n})
	}
	return deps
}

func (dlp *depLockProjector) parseName(text string) string {
	matches := projectNameRegex.FindStringSubmatch(text)
	return strings.TrimSpace(matches[len(matches)-1])
}
