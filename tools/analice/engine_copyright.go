package main

import (
	"regexp"
	"strings"
)

type copyrightEngine struct {
	rootCopyrightHolders   []string
	depsVsCopyrightHolders map[dependency][]string
	curDep                 dependency
	curCopyrightHolders    []string
}

func (c *copyrightEngine) initialize() {
	c.depsVsCopyrightHolders = make(map[dependency][]string)
}

func (c *copyrightEngine) start(p project, d dependency) error {
	c.curDep = d
	c.curCopyrightHolders = make([]string, 0)
	return nil
}

func (c *copyrightEngine) stop(p project, d dependency) error {
	holders := make([]string, len(c.curCopyrightHolders))
	copy(holders, c.curCopyrightHolders)
	if d.name == "" {
		c.rootCopyrightHolders = holders
	} else {
		c.depsVsCopyrightHolders[d] = holders
	}
	return nil
}

var copyrightRegex = regexp.MustCompile(`^>*\s*[Cc]opyright\s(\([Cc]\)|©)*(?P<CLINE>.*)`)

var copyrightNoAuthorRegex = regexp.MustCompile(`([Cc]opyright notice|[Cc]opyright license to reproduce)`)

func (c *copyrightEngine) analyze(line string) error {
	if copyrightRegex.MatchString(line) && !copyrightNoAuthorRegex.MatchString(line) {
		c.curCopyrightHolders = append(c.curCopyrightHolders, c.parseCopyrightHolder(line))
	}
	return nil
}

func (c *copyrightEngine) parseCopyrightHolder(text string) string {
	matches := copyrightRegex.FindStringSubmatch(text)
	return strings.TrimSpace(matches[len(matches)-1])
}
