package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

var stdlibPkgs = loadStdlib()

type xPackagesScanner struct {
}

func (x xPackagesScanner) scan(opt depScannerOptions, patterns ...string) ([]pkg, error) {
	var ps []pkg
	for _, perm := range opt.envPermutations {
		var osVariant []string
		osVariant = append(osVariant, os.Environ()...)
		osVariant = append(osVariant, perm...)
		scanResults, err := x.cscan(packages.Config{
			Mode:  packages.LoadImports,
			Tests: opt.tests,
			Env:   osVariant,
		}, patterns...)
		if err != nil {
			return nil, err
		}
		ps = x.merge(ps, scanResults)
	}
	return ps, nil
}

func (x xPackagesScanner) merge(ps []pkg, others []pkg) []pkg {
	for _, scanResult := range others {
		if !x.isContained(scanResult, ps) {
			ps = append(ps, scanResult)
		}
	}
	return ps
}

func (x xPackagesScanner) isContained(p pkg, others []pkg) bool {
	for _, other := range others {
		if p.PkgPath == other.PkgPath {
			return true
		}
	}
	return false
}

func (x xPackagesScanner) cscan(cfg packages.Config, patterns ...string) ([]pkg, error) {
	ps, err := packages.Load(&cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("unable to load '%s': %v", patterns, err)
	}
	filtered := x.noStdlib(ps)
	condensed := x.truncate(filtered)
	pkgs := x.markRoots(condensed, patterns)
	return pkgs, nil
}

func (x xPackagesScanner) markRoots(ps []*packages.Package, patterns []string) []pkg {
	pkgs := make([]pkg, 0, len(ps))
	for _, p := range ps {
		pkgs = append(pkgs, pkg{Package: p, root: x.isRoot(p, patterns)})
	}
	return pkgs
}

func (x xPackagesScanner) noStdlib(ps []*packages.Package) []*packages.Package {
	var filtered []*packages.Package
	packages.Visit(ps, func(p *packages.Package) bool {
		_, ok := stdlibPkgs[p.PkgPath]
		if ok {
			return false
		}
		filtered = append(filtered, p)
		return true
	}, func(p *packages.Package) {})
	return filtered
}

func (x xPackagesScanner) truncate(ps []*packages.Package) []*packages.Package {
	var roots []*packages.Package
	for i, p := range ps {
		hasParent := x.hasParent(p, ps[:i]) || x.hasParent(p, ps[i+1:])
		if !hasParent {
			roots = append(roots, p)
		}
	}
	return roots
}

func (x xPackagesScanner) hasParent(p *packages.Package, others []*packages.Package) bool {
	for _, q := range others {
		if strings.HasPrefix(p.PkgPath, q.PkgPath) {
			return true
		}
	}
	return false
}

func (x xPackagesScanner) isRoot(p *packages.Package, patterns []string) bool {
	for _, pattern := range patterns {
		if pattern == p.PkgPath {
			return true
		}
	}
	return false
}
