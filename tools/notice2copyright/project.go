package main

import (
	"net/url"
	"path/filepath"
)

// codebeat:disable[TOO_MANY_IVARS]

type project struct {
	name             string
	url              string
	license          license
	copyrightHolders []string
	isRoot           bool
}

// codebeat:enable[TOO_MANY_IVARS]

func (p *project) dir() string {
	if p.isRoot {
		return "."
	}
	u, err := url.Parse(p.url)
	assertOrFatal(err == nil, "unable to parse URL for project '%s' (%s): %v", p.name, p.url, err)
	return filepath.Join("vendor", u.Host, filepath.FromSlash(u.Path))
}
