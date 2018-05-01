package main

var licenseFiles = []string{"LICENSE", "LICENSE.md", "LICENSE.txt", "LICENSE.rst", "COPYING", "License", "MIT-LICENSE.txt"}

type license struct {
	name      string
	shortName string
	text      string
}

var licenses []license
