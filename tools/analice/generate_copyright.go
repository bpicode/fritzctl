package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var copyrightCmd = &cobra.Command{
	Use:  "copyright /path/to/project",
	RunE: genCopyright,
}

func init() {
	generateCmd.AddCommand(copyrightCmd)
}

func genCopyright(_ *cobra.Command, args []string) error {
	dir := projectDir(args)
	p, err := getProjector().project(dir)
	if err != nil {
		return fmt.Errorf("failed to get project dependencies: %v", err)
	}
	genCopyrightOf(p)
	return nil
}

func genCopyrightOf(p *project) {
	bowEngine := bowLicenseEngine{}
	cprEngine := copyrightEngine{}
	a := newAnalyzer(&bowEngine, &cprEngine)
	a.run(*p)
	c := cprWriter(*p, bowEngine, cprEngine)
	c.writeTo(os.Stdout)
}

func cprWriter(p project, b bowLicenseEngine, c copyrightEngine) copyrightWriter {
	ds := joinLicAndCph(b, c)
	w := copyrightWriter{
		project:              p,
		rootCopyrightHolders: c.rootCopyrightHolders,
		rootLic:              b.rootLic,
		deps:                 ds,
	}
	return w
}

func joinLicAndCph(b bowLicenseEngine, c copyrightEngine) []copyrightedLicenseDep {
	depsVsLic := transposeLicDep(b.licsVsDeps)
	var ds []copyrightedLicenseDep
	for d, l := range depsVsLic {
		cld := copyrightedLicenseDep{}
		cld.dep = d
		cld.lic = l
		cld.cph = c.depsVsCopyrightHolders[d]
		ds = append(ds, cld)
	}
	return ds
}

func transposeLicDep(m map[license][]dependency) map[dependency]license {
	t := make(map[dependency]license)
	for l, ds := range m {
		for _, d := range ds {
			t[d] = l
		}
	}
	return t
}
