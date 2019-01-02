package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

var copyrightRegex = regexp.MustCompile(`^>*\s*[Cc]opyright\s(\([Cc]\)|©)*(?P<CLINE>.*)`)

var copyrightNoAuthorRegex = regexp.MustCompile(`([Cc]opyright notice|[Cc]opyright license to reproduce)`)

type licenseFileAnalizer struct {
}

func (l licenseFileAnalizer) analyze(ps ...pkgWithLicenseFile) ([]licPkg, []error) {
	lps := make([]licPkg, len(ps))
	errs := make([]error, len(ps))
	bowEngine, cprEngine := l.engines()
	for i, p := range ps {
		ptmp := p
		lps[i] = licPkg{pkgWithLicenseFile: &ptmp}
		lps[i].license, lps[i].copyrightHolders, errs[i] = l.analyzeOne(p, bowEngine, cprEngine)
	}
	return lps, errs
}

func (l licenseFileAnalizer) analyzeOne(p pkgWithLicenseFile, bowEngine bowEngine, cprEngine copyrightRegexHeuristic) (*license, []string, error) {
	file, err := os.Open(p.licFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to open license file '%s': %v", p.licFilePath, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	bowEngine.start()
	cprEngine.start()
	for scanner.Scan() {
		line := scanner.Text()
		bowEngine.analyze(line)
		cprEngine.analyze(line)
	}

	err = bowEngine.stop()
	cprEngine.stop()
	return &bowEngine.bestGuess, cprEngine.holders, err
}

func (l licenseFileAnalizer) engines() (bowEngine, copyrightRegexHeuristic) {
	bowEngine := bowEngine{}
	cprEngine := copyrightRegexHeuristic{}
	bowEngine.initialize()
	cprEngine.initialize()
	return bowEngine, cprEngine
}
