package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
)

type copyrightRenderer struct {
}

func (c copyrightRenderer) render(w io.Writer, ps ...licPkg) error {
	if len(ps) == 0 {
		return errors.New("cannot render copyright, no packages given")
	}
	c.sort(ps)
	c.writeHead(w, ps[0])
	c.writeRootProjects(w, ps...)
	c.write3rdParties(w, ps...)
	c.writeLics(w, ps...)
	return nil
}

func (c copyrightRenderer) sort(ps []licPkg) {
	sort.SliceStable(ps, func(i, j int) bool {
		iRoot := ps[i].root
		jRoot := ps[j].root
		if iRoot && !jRoot {
			return true
		}
		if !iRoot && jRoot {
			return false
		}
		return ps[i].PkgPath < ps[j].PkgPath
	})
}

func (c copyrightRenderer) writeHead(w io.Writer, root licPkg) {
	fmt.Fprintln(w, "Format: https://www.debian.org/doc/packaging-manuals/copyright-format/1.0/")
	fmt.Fprintf(w, "Upstream-Name: %s\n", root.short())
	fmt.Fprintf(w, "Source: %s\n", root.url())
	fmt.Fprintln(w, "")
}

func (c copyrightRenderer) writeRootProjects(w io.Writer, ps ...licPkg) {
	c.walk(w, func(p licPkg) bool {
		return c.isRoot(p) && c.hasLicense(p)
	}, c.writeRootProject, ps...)
}

func (c copyrightRenderer) writeRootProject(w io.Writer, root licPkg) {
	fmt.Fprintln(w, "Files: *")
	c.writeCph(w, root.copyrightHolders)
	fmt.Fprintf(w, "License: %s\n\n", root.license.shortName)
}

type pkgFilterFunc func(licPkg) bool
type pkgConsumerFunc func(io.Writer, licPkg)

func (c copyrightRenderer) walk(w io.Writer, filter pkgFilterFunc, consumer pkgConsumerFunc, ps ...licPkg) {
	for _, p := range ps {
		if !filter(p) {
			continue
		}
		consumer(w, p)
	}
}

func (c copyrightRenderer) isRoot(p licPkg) bool {
	return p.root
}

func (c copyrightRenderer) isNotRoot(p licPkg) bool {
	return !p.root
}

func (c *copyrightRenderer) hasLicense(p licPkg) bool {
	if p.license == nil {
		log.Printf("no license for '%s', skipping", p.PkgPath)
		return false
	}
	return true
}

func (c copyrightRenderer) writeCph(w io.Writer, cph []string) {
	if len(cph) == 0 {
		fmt.Fprintf(w, "Copyright: %s\n", "unknown")
		return
	}
	fmt.Fprintf(w, "Copyright: %s\n", cph[0])
	for _, other := range cph[1:] {
		fmt.Fprintf(w, "           %s\n", other)
	}
}

func (c copyrightRenderer) write3rdParties(w io.Writer, ps ...licPkg) {
	c.walk(w, func(p licPkg) bool {
		return c.isNotRoot(p) && c.hasLicense(p)
	}, c.write3rdParty, ps...)

}

func (c copyrightRenderer) write3rdParty(w io.Writer, p licPkg) {
	fmt.Fprintf(w, "Files: vendor/%s/*\n", p.PkgPath)
	c.writeCph(w, p.copyrightHolders)
	fmt.Fprintf(w, "License: %s\n\n", p.license.shortName)
}

func (c copyrightRenderer) writeLics(w io.Writer, ps ...licPkg) {
	m := make(map[license]bool)
	c.walk(w, c.hasLicense, func(w io.Writer, p licPkg) {
		m[*p.license] = true
	}, ps...)

	allLics := c.licensesSorted(m)
	for _, l := range allLics {
		c.writeLic(w, l)
	}
}

func (c copyrightRenderer) licensesSorted(m map[license]bool) []license {
	allLics := make([]license, 0)
	for l := range m {
		allLics = append(allLics, l)
	}
	sort.Slice(allLics, func(i, j int) bool {
		return allLics[i].shortName < allLics[j].shortName
	})
	return allLics
}

func (c copyrightRenderer) writeLic(w io.Writer, l license) {
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
