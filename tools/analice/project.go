package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
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

type goModProjector struct {
}

func (gmp *goModProjector) project(dir string) (*project, error) {
	join := path.Join(dir, "go.mod")
	file, err := os.Open(join)
	if err != nil {
		return nil, fmt.Errorf("failed to open go.sum: %v", err)
	}
	defer file.Close()
	return &project{root: dir, dependencies: gmp.parse(file)}, nil
}

func (gmp *goModProjector) parse(r io.Reader) []dependency {
	s := bufio.NewScanner(r)
	namesVsDeps := make(map[string]dependency)
	withinRequireBlock := false
	for s.Scan() {
		text := s.Text()
		if text == "require (" {
			withinRequireBlock = true
			continue
		}
		if text == ")" {
			withinRequireBlock = false
			continue
		}
		if withinRequireBlock {
			gmp.condAppendProject(namesVsDeps, text)
		}
	}
	return gmp.values(namesVsDeps)
}

func (gmp *goModProjector) condAppendProject(m map[string]dependency, text string) {
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return
	}
	name := fields[0]
	m[name] = dependency{name: name}
}

func (gmp *goModProjector) values(m map[string]dependency) []dependency {
	ds := make([]dependency, 0, len(m))
	for _, v := range m {
		ds = append(ds, v)
	}
	return ds
}
