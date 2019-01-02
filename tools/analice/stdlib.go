package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

var sym = regexp.MustCompile(`^pkg (\S+).*?, (?:var|func|type|const) ([a-zA-Z]\w*)`)

var apiFiles = []string{
	"go1.txt",
	"go1.1.txt",
	"go1.2.txt",
	"go1.3.txt",
	"go1.4.txt",
	"go1.5.txt",
	"go1.6.txt",
	"go1.7.txt",
	"go1.8.txt",
	"go1.9.txt",
	"go1.10.txt",
	"go1.11.txt",
}

type stdlib map[string]bool

func loadStdlib() stdlib {
	s := stdlib{}
	files := s.openAPIFiles()
	defer s.closeAll(files)
	s.putStdlibImports(files)
	s.putMissedStdlibImports()
	return s
}

func (s stdlib) openAPIFiles() []*os.File {
	var files []*os.File
	for _, apiFile := range apiFiles {
		files = append(files, s.open(apiFile))
	}
	return files
}

func (s stdlib) open(name string) *os.File {
	f, _ := os.Open(filepath.Join(runtime.GOROOT(), "api", name))
	return f
}

func (s stdlib) closeAll(files []*os.File) {
	for _, f := range files {
		if f != nil {
			f.Close()
		}
	}
}

func (s stdlib) putStdlibImports(files []*os.File) {
	var readers []io.Reader
	for _, f := range files {
		readers = append(readers, f)
	}
	mr := io.MultiReader(readers...)
	sc := bufio.NewScanner(mr)
	for sc.Scan() {
		l := sc.Text()
		s.putLine(l)
	}
}

func (s stdlib) putLine(l string) {
	if m := sym.FindStringSubmatch(l); m != nil {
		path := m[1]
		if _, ok := s[path]; !ok {
			s[path] = true
		}
	}
}

func (s stdlib) putMissedStdlibImports() {
	s["unsafe"] = true
	s["testing/internal/testdeps"] = true
	s["internal/testlog"] = true
}
