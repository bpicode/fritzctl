package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

type copyrightWriter struct {
	root project
	deps []project
}

func (c *copyrightWriter) writeTo(w io.Writer) {
	c.writeHead(w)
	c.writeProject(w, c.root)
	c.writeDeps(w)
	c.writeLics(w)
}

func (c *copyrightWriter) writeHead(w io.Writer) {
	fmt.Fprintln(w, "Format: https://www.debian.org/doc/packaging-manuals/copyright-format/1.0/")
	fmt.Fprintf(w, "Upstream-Name: %s\n", c.root.name)
	fmt.Fprintf(w, "Source: %s\n", c.root.url)
	fmt.Fprintln(w, "")
}

func (c *copyrightWriter) writeProject(w io.Writer, p project) {
	if p.isRoot {
		fmt.Fprintln(w, "Files: *")
	} else {
		fmt.Fprintf(w, "Files: %s/*\n", p.dir())
	}

	fmt.Fprintf(w, "Copyright: %s\n", p.copyrightHolders[0])
	for _, other := range p.copyrightHolders[1:] {
		fmt.Fprintf(w, "           %s\n", other)
	}
	fmt.Fprintf(w, "License: %s\n\n", p.license.shortName)
}

func (c *copyrightWriter) writeDeps(w io.Writer) {
	sort.Slice(c.deps, func(i, j int) bool {
		return c.deps[i].name < c.deps[j].name
	})
	for _, d := range c.deps {
		c.writeProject(w, d)
	}
}

func (c *copyrightWriter) writeLics(w io.Writer) {
	m := c.licSet()
	allLics := c.sortByShortName(m)
	for _, l := range allLics {
		c.writeLic(w, l)
	}
}

func (c *copyrightWriter) licSet() map[license]interface{} {
	m := make(map[license]interface{})
	m[c.root.license] = nil
	for _, p := range c.deps {
		m[p.license] = nil
	}
	return m
}

func (c *copyrightWriter) sortByShortName(m map[license]interface{}) []license {
	allLics := make([]license, 0)
	for l := range m {
		allLics = append(allLics, l)
	}
	sort.Slice(allLics, func(i, j int) bool {
		return allLics[i].shortName < allLics[j].shortName
	})
	return allLics
}

func (c *copyrightWriter) writeLic(w io.Writer, l license) {
	fmt.Printf("License: %s\n", l.shortName)
	s := bufio.NewScanner(bytes.NewBufferString(l.text))
	for s.Scan() {
		line := s.Text()
		if strings.TrimSpace(line) == "" {
			fmt.Fprintln(w, " .")
		} else {
			fmt.Fprintf(w, " %s\n", line)
		}
	}
	fmt.Fprintln(w, "")
}
