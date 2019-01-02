package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/packages"
)

func TestResolve(t *testing.T) {
	scanner := xPackagesScanner{}
	pkgs, err := scanner.scan(depScannerOptions{
		envPermutations: [][]string{nil},
		tests:           false,
	}, "github.com/bpicode/fritzctl")
	assert.NoError(t, err)
	assert.NotEmpty(t, pkgs)

	r := licenceFileResolver{}
	withLicenseFiles, errs := r.resolve(pkgs...)
	assert.NotEmpty(t, withLicenseFiles)
	for _, err = range errs {
		assert.NoError(t, err)
	}

	for i, wlf := range withLicenseFiles {
		assert.Equal(t, pkgs[i].PkgPath, wlf.PkgPath)
	}

	for _, wlf := range withLicenseFiles {
		assert.NotEmpty(t, wlf.licFilePath)
	}
}

func TestResolveOne(t *testing.T) {
	r := licenceFileResolver{}
	_, err := r.resolveOne(pkg{Package: &packages.Package{
		PkgPath: "lol.com/hypthetical/package",
	}})
	assert.Error(t, err)

	_, err = r.resolveOne(pkg{Package: &packages.Package{
		PkgPath: "lol.com/hypthetical/package",
		GoFiles: []string{"/path/to/some/file.go"},
	}})
	assert.Error(t, err)
}
