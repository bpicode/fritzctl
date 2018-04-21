package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_licenseGuesser_openLicense(t *testing.T) {
	a := analyzer{}
	_, err := a.openLicense("/that/should/not/work/ever")
	assert.Error(t, err)
}

func Test_licenseGuesser_guess_file_not_found(t *testing.T) {
	a := newAnalyzer(&bowLicenseEngine{})
	p := project{
		root:         "/points/to/nowhere",
		dependencies: []dependency{{name: "github.com/acme/something"}},
	}
	assert.NotPanics(t, func() { a.run(p) })
}

func Test_licenseGuesser_guess_one_undetermined(t *testing.T) {
	pRoot, err := ioutil.TempDir("", "analice-test")
	assert.NoError(t, err)
	defer os.RemoveAll(pRoot)
	err = os.MkdirAll(path.Join(pRoot, "vendor", "dependency"), 0755)
	assert.NoError(t, err)
	err = ioutil.WriteFile(path.Join(pRoot, "vendor", "dependency", "LICENSE"), nil, 0644)
	assert.NoError(t, err)

	a := newAnalyzer(&bowLicenseEngine{})
	assert.NoError(t, err)
	p := project{
		root:         pRoot,
		dependencies: []dependency{{name: "dependency"}},
	}
	assert.NotPanics(t, func() { a.run(p) })
}
