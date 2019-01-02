package main

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/tools/go/packages"
)

type depScanner interface {
	scan(opt depScannerOptions, patterns ...string) ([]pkg, error)
}

type depScannerOptions struct {
	tests           bool
	envPermutations [][]string
}

type pkg struct {
	*packages.Package
	root bool
}

func (p *pkg) short() string {
	split := strings.Split(p.PkgPath, "/")
	padded := append([]string{""}, split...)
	return padded[len(padded)-1]
}

func (p *pkg) url() string {
	return fmt.Sprintf("https://%s", p.PkgPath)
}

type resolver interface {
	resolve(ps ...pkg) ([]pkgWithLicenseFile, []error)
}

type pkgWithLicenseFile struct {
	*pkg
	licFilePath string
}

type licenseAnalyzer interface {
	analyze(ps ...pkgWithLicenseFile) ([]licPkg, []error)
}

type licPkg struct {
	*pkgWithLicenseFile
	copyrightHolders []string
	license          *license
}

type license struct {
	name      string
	shortName string
	text      string
}

type renderer interface {
	render(w io.Writer, ps ...licPkg) error
}

func newDepScanner() depScanner {
	return xPackagesScanner{}
}

func newResolver() resolver {
	return licenceFileResolver{}
}

func newAnalyzer() licenseAnalyzer {
	return licenseFileAnalizer{}
}

func newDebianCopyrightRenderer() renderer {
	return copyrightRenderer{}
}

func newNoticeFileRenderer() renderer {
	return noticeRenderer{}
}
