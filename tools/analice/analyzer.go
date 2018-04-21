package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

type analyzer struct {
	lineAnalyzers []lineAnalyzer
}

type lineAnalyzer interface {
	initialize()
	start(p project, d dependency) error
	analyze(line string) error
	stop(p project, d dependency) error
}

func (l *analyzer) register(a lineAnalyzer) {
	l.lineAnalyzers = append(l.lineAnalyzers, a)
}

func (l *analyzer) run(p project) {
	l.initialize()
	l.analyzeRoot(p)
	l.analyzeDeps(p)
}

func (l *analyzer) initialize() {
	for _, a := range l.lineAnalyzers {
		a.initialize()
	}
}

func (l *analyzer) analyzeRoot(p project) {
	empty := dependency{}
	l.start(p, empty)
	l.analyzeDir(p.root, p, empty)
	l.stop(p, empty)
}

func (l *analyzer) start(p project, d dependency) {
	for _, a := range l.lineAnalyzers {
		a.start(p, d)
	}
}

func (l *analyzer) stop(p project, d dependency) {
	for _, a := range l.lineAnalyzers {
		err := a.stop(p, d)
		if err != nil {
			log.Printf("problem analyzing license for '%s': %v", d.name, err)
		}
	}
}

func (l *analyzer) analyzeDeps(p project) {
	for _, d := range p.dependencies {
		l.start(p, d)
		err := l.analyzeDep(p, d)
		if err != nil {
			log.Printf("problem analyzing license for '%s': %v", d.name, err)
			continue
		}
		l.stop(p, d)
	}
}

func (l *analyzer) analyzeDep(p project, d dependency) error {
	dir := path.Join(p.root, "vendor", d.name)
	return l.analyzeDir(dir, p, d)
}

func (l *analyzer) analyzeDir(dir string, p project, d dependency) error {
	file, err := l.openLicense(dir)
	if err != nil {
		return fmt.Errorf("unable to open license file in directory '%s': %v", dir, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, a := range l.lineAnalyzers {
			a.analyze(line)
		}
	}
	return nil
}

func (l *analyzer) openLicense(dir string) (*os.File, error) {
	for _, name := range licenseFiles {
		f, err := os.Open(filepath.Join(dir, name))
		if err == nil {
			return f, nil
		}
	}
	return nil, fmt.Errorf("no license file found in directory '%s' (have tried %v)", dir, licenseFiles)
}
