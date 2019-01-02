package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXPackagesScanner(t *testing.T) {
	scanner := xPackagesScanner{}
	pkgs, err := scanner.scan(depScannerOptions{
		envPermutations: [][]string{nil, {"GOOS=windows"}},
		tests:           false,
	}, "github.com/bpicode/fritzctl")
	assert.NoError(t, err)
	assert.NotEmpty(t, pkgs)
}

func TestXPackagesScannerErr(t *testing.T) {
	scanner := xPackagesScanner{}
	_, err := scanner.scan(depScannerOptions{
		envPermutations: [][]string{{"GOOS=this_should_not_work"}},
		tests:           false,
	}, "some/package")
	assert.Error(t, err)
}
