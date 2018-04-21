package main

import (
	"fmt"
	"io"
	"sort"
)

type noticeWriter struct {
}

type section struct {
	license      string
	dependencies []dependency
}

func (n *noticeWriter) writeTo(w io.Writer, m map[string][]dependency) {
	n.writeIntroWords(w)
	secs := n.sort(m)
	writeLicenseBlocks(w, secs)
}

func (n *noticeWriter) writeIntroWords(w io.Writer) {
	fmt.Fprintf(w, "%s\n\n", "# 3rd-Party Software License Notice")
	fmt.Fprintf(w, "%s\n\n", "This software includes the following software and licenses:")
}

func (n *noticeWriter) sort(m map[string][]dependency) []section {
	var secs []section
	for k, v := range m {
		sec := n.ordered(k, v)
		secs = append(secs, sec)
	}
	sort.Slice(secs, func(i, j int) bool {
		return secs[i].license < secs[j].license
	})
	return secs
}

func (n *noticeWriter) ordered(lic string, deps []dependency) section {
	sec := section{license: lic, dependencies: make([]dependency, len(deps))}
	copy(sec.dependencies, deps)
	sort.Slice(sec.dependencies, func(i, j int) bool {
		return sec.dependencies[i].short() < sec.dependencies[j].short()
	})
	return sec
}

func writeLicenseBlocks(w io.Writer, secs []section) {
	for _, s := range secs {
		fmt.Fprintf(w, "%s\n", "========================================================================")
		fmt.Fprintf(w, "%s\n", s.license)
		fmt.Fprintf(w, "%s\n\n", "========================================================================")
		fmt.Fprintf(w, "%s\n\n", "The following software have components provided under the terms of this license:")

		for _, d := range s.dependencies {
			fmt.Fprintf(w, "- %s (from https://%s)\n", d.short(), d.name)
		}
		fmt.Fprintln(w)
	}
}
