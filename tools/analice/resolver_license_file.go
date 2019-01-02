package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var licenseFiles = []string{"LICENSE", "LICENSE.md", "LICENSE.txt", "LICENSE.rst", "COPYING", "License", "MIT-LICENSE.txt"}

type licenceFileResolver struct {
}

func (r licenceFileResolver) resolve(ps ...pkg) ([]pkgWithLicenseFile, []error) {
	pswlf := make([]pkgWithLicenseFile, len(ps))
	errs := make([]error, len(ps))
	for i, p := range ps {
		ptmp := p
		withLic := pkgWithLicenseFile{}
		withLic.pkg = &ptmp
		withLic.licFilePath, errs[i] = r.resolveOne(p)
		pswlf[i] = withLic
	}
	return pswlf, errs
}

func (r licenceFileResolver) resolveOne(p pkg) (string, error) {
	if len(p.GoFiles) == 0 {
		return "", fmt.Errorf("no go files for package '%v', don't know where to search for license", p.PkgPath)
	}
	startingDir := filepath.Dir(p.GoFiles[0])

	for dir := startingDir; r.stillInScope(dir, p.PkgPath); dir = filepath.Dir(dir) {
		path := r.searchLicenseFile(dir)
		if path != "" {
			return path, nil
		}
	}

	return "", fmt.Errorf("no license file found, starting from '%s' (have tried %v)", p.GoFiles[0], licenseFiles)
}

func (r licenceFileResolver) stillInScope(dir string, pkgPath string) bool {
	gp := strings.Index(dir, os.Getenv("GOPATH"))
	dgp := strings.Index(dir, build.Default.GOPATH)
	gr := strings.Index(dir, runtime.GOROOT())
	if gp >= 0 || dgp >= 0 || gr >= 0 {
		return true
	}
	return strings.Contains(dir, pkgPath)
}

func (r licenceFileResolver) searchLicenseFile(dir string) string {
	for _, name := range licenseFiles {
		path := filepath.Join(dir, name)
		_, err := os.Stat(path)
		if err == nil {
			return path
		}
	}
	return ""
}
