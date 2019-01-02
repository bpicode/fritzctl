package main

import (
	"fmt"
	"io"
	"log"
	"sort"
)

type noticeRenderer struct {
}

type noticeGroup struct {
	license string
	modules []licPkg
}

func (n noticeRenderer) render(w io.Writer, ps ...licPkg) error {
	n.writeIntroWords(w)
	m := n.groupByLicense(ps...)
	gs := n.groups(m)
	n.writeLicenseBlocks(w, gs)
	return nil
}

func (n noticeRenderer) writeIntroWords(w io.Writer) {
	fmt.Fprintf(w, "%s\n\n", "# 3rd-Party Software License Notice")
	fmt.Fprintf(w, "%s\n\n", "This software includes the following software and licenses:")
}

func (n noticeRenderer) groups(m map[string][]licPkg) []noticeGroup {
	var groups []noticeGroup
	for k, v := range m {
		sec := n.ordered(k, v)
		groups = append(groups, sec)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].license < groups[j].license
	})
	return groups
}

func (n noticeRenderer) ordered(lic string, ps []licPkg) noticeGroup {
	g := noticeGroup{license: lic, modules: make([]licPkg, len(ps))}
	copy(g.modules, ps)
	sort.Slice(g.modules, func(i, j int) bool {
		return g.modules[i].short() < g.modules[j].short()
	})
	return g
}

func (n noticeRenderer) writeLicenseBlocks(w io.Writer, gs []noticeGroup) {
	for _, g := range gs {
		fmt.Fprintf(w, "%s\n", "========================================================================")
		fmt.Fprintf(w, "%s\n", g.license)
		fmt.Fprintf(w, "%s\n\n", "========================================================================")
		fmt.Fprintf(w, "%s\n\n", "The following software have components provided under the terms of this license:")

		for _, d := range g.modules {
			fmt.Fprintf(w, "- %s (from %s)\n", d.short(), d.url())
		}
		fmt.Fprintln(w)
	}
}

func (n noticeRenderer) groupByLicense(ps ...licPkg) map[string][]licPkg {
	m := make(map[string][]licPkg)
	for _, p := range ps {
		if p.root {
			continue
		}
		if p.license == nil {
			log.Printf("no license for '%s', skipping", p.url())
			continue
		}
		m[p.license.name] = append(m[p.license.name], p)
	}
	return m
}
