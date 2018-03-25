package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

var projectRegex = regexp.MustCompile(`\s*-*\s*(?P<NAME>\S+)\s*\(\s*from\s*(?P<URL>\S+)\).*`)

func parseNotice(name string) []project {
	s := readNotice(name)
	projects := collectProjects(bytes.NewBuffer(s))
	return projects
}

func readNotice(name string) []byte {
	f, err := os.Open(name)
	assertOrFatal(err == nil, "failed to open '%s': %v", name, err)
	defer f.Close()
	s, err := ioutil.ReadAll(f)
	assertOrFatal(err == nil, "unable to read '%s: %v", name, err)
	return s
}

func collectProjects(r io.Reader) []project {
	s := bufio.NewScanner(r)
	var projects []project
	var curLic license
	for s.Scan() {
		text := s.Text()
		updateCurrentLicense(text, &curLic)
		projects = addProject(projects, text, curLic)
	}
	return projects
}

func addProject(projects []project, text string, lic license) []project {
	if tokens := projectRegex.FindStringSubmatch(text); len(tokens) == 3 {
		projects = append(projects, project{license: lic, name: tokens[1], url: tokens[2]})
	}
	return projects
}

func updateCurrentLicense(text string, current *license) {
	if lic, ok := licenses[text]; ok {
		*current = lic
	}
}
