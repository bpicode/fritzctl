package main

import "strings"

type copyrightRegexHeuristic struct {
	holders []string
}

func (c *copyrightRegexHeuristic) initialize() {
}

func (c *copyrightRegexHeuristic) start() {
	c.holders = make([]string, 0)
}

func (c *copyrightRegexHeuristic) stop() error {
	holders := make([]string, len(c.holders))
	copy(holders, c.holders)
	return nil
}

func (c *copyrightRegexHeuristic) analyze(line string) error {
	if copyrightRegex.MatchString(line) && !copyrightNoAuthorRegex.MatchString(line) {
		c.holders = append(c.holders, c.parseCopyrightHolder(line))
	}
	return nil
}

func (c *copyrightRegexHeuristic) parseCopyrightHolder(text string) string {
	matches := copyrightRegex.FindStringSubmatch(text)
	return strings.TrimSpace(matches[len(matches)-1])
}
