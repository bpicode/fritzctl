package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// codebeat:disable[TOO_MANY_IVARS]
type copyrightWriter struct {
	project              project
	rootLic              license
	rootCopyrightHolders []string
	deps                 []copyrightedLicenseDep
}

// codebeat:enable[TOO_MANY_IVARS]

type copyrightedLicenseDep struct {
	dep dependency
	cph []string
	lic license
}

func (c *copyrightWriter) writeTo(w io.Writer) {
	c.writeHead(w)
	c.writeRootProject(w)
	c.writeDeps(w)
	c.writeLics(w)
}

func (c *copyrightWriter) writeHead(w io.Writer) {
	fmt.Fprintln(w, "Format: https://www.debian.org/doc/packaging-manuals/copyright-format/1.0/")
	fmt.Fprintf(w, "Upstream-Name: %s\n", c.project.short())
	fmt.Fprintf(w, "Source: %s\n", c.project.url())
	fmt.Fprintln(w, "")
}

func (c *copyrightWriter) writeRootProject(w io.Writer) {
	fmt.Fprintln(w, "Files: *")
	c.writeCph(w, c.rootCopyrightHolders)
	fmt.Fprintf(w, "License: %s\n\n", c.rootLic.shortName)
}

func (c *copyrightWriter) writeDeps(w io.Writer) {
	sort.Slice(c.deps, func(i, j int) bool {
		return c.deps[i].dep.name < c.deps[j].dep.name
	})
	for _, d := range c.deps {
		c.writeDep(w, d)
	}
}

func (c *copyrightWriter) writeDep(w io.Writer, d copyrightedLicenseDep) {
	fmt.Fprintf(w, "Files: vendor/%s/*\n", d.dep.name)
	c.writeCph(w, d.cph)
	fmt.Fprintf(w, "License: %s\n\n", d.lic.shortName)
}

func (c *copyrightWriter) writeCph(w io.Writer, cph []string) {
	if len(cph) == 0 {
		fmt.Fprintf(w, "Copyright: %s\n", "unknown")
		return
	}
	fmt.Fprintf(w, "Copyright: %s\n", cph[0])
	for _, other := range cph[1:] {
		fmt.Fprintf(w, "           %s\n", other)
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
	m[c.rootLic] = nil
	for _, d := range c.deps {
		m[d.lic] = nil
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
