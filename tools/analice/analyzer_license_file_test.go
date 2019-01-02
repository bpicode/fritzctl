package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyze(t *testing.T) {
	pkgs, _ := xPackagesScanner{}.scan(depScannerOptions{
		envPermutations: [][]string{nil},
		tests:           false,
	}, "github.com/bpicode/fritzctl")
	withLicenseFiles, _ := licenceFileResolver{}.resolve(pkgs...)

	licPkgs, errs := licenseFileAnalizer{}.analyze(withLicenseFiles...)
	for _, err := range errs {
		assert.NoError(t, err)
	}
	assert.NotEmpty(t, licPkgs)
}

func TestAnalyzeOne(t *testing.T) {
	a := licenseFileAnalizer{}
	e1, e2 := a.engines()
	_, _, err := a.analyzeOne(pkgWithLicenseFile{licFilePath: ""}, e1, e2)
	assert.Error(t, err)
}
