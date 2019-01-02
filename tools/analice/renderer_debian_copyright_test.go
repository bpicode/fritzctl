package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/packages"
)

func TestRender(t *testing.T) {
	pkgs, _ := xPackagesScanner{}.scan(depScannerOptions{
		envPermutations: [][]string{nil},
		tests:           false,
	}, "github.com/bpicode/fritzctl")
	withLicenseFiles, _ := licenceFileResolver{}.resolve(pkgs...)

	licPkgs, _ := licenseFileAnalizer{}.analyze(withLicenseFiles...)
	copyrightRenderer{}.render(os.Stdout, licPkgs...)
}

func Test_copyrightRenderer_sort(t *testing.T) {
	r := copyrightRenderer{}
	ps := []licPkg{somePackage("x/y", false, &apache2Short), somePackage("b/c", true, &mit), somePackage("a/b", false, &bsd2)}
	r.sort(ps)
	assert.Equal(t, "b/c", ps[0].PkgPath)
	assert.Equal(t, "a/b", ps[1].PkgPath)
	assert.Equal(t, "x/y", ps[2].PkgPath)
}

func somePackage(path string, root bool, l *license) licPkg {
	return licPkg{
		license: l,
		pkgWithLicenseFile: &pkgWithLicenseFile{
			pkg: &pkg{
				root: root,
				Package: &packages.Package{
					PkgPath: path,
				},
			},
		},
	}
}

func Test_copyrightRenderer_render_no_pkgs(t *testing.T) {
	c := copyrightRenderer{}
	err := c.render(ioutil.Discard)
	assert.Error(t, err)
}

func Test_copyrightRenderer_hasLicense(t *testing.T) {
	c := copyrightRenderer{}
	assert.False(t, c.hasLicense(somePackage("x", false, nil)))
	assert.True(t, c.hasLicense(somePackage("y", false, &mit)))
}
