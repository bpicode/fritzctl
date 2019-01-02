package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type genericRenderCmd struct {
	scanner  depScanner
	renderer renderer
}

func (grc genericRenderCmd) run(cmd *cobra.Command, args []string) error {
	ps, err := grc.runScan(cmd, args)
	if err != nil {
		return fmt.Errorf("failed to obtain dependencies: %sv", err)
	}
	err = grc.runRender(ps)
	return err
}

func (grc genericRenderCmd) runScan(cmd *cobra.Command, args []string) ([]pkg, error) {
	includeTests, _ := cmd.Flags().GetBool("tests")
	ps, err := grc.scanner.scan(depScannerOptions{
		tests:           includeTests,
		envPermutations: grc.envPerm(cmd),
	}, args...)
	return ps, err
}

func (grc genericRenderCmd) runRender(ps []pkg) error {
	withLicenseFiles, errs := newResolver().resolve(ps...)
	grc.logErrs(errs)
	licPkgs, errs := newAnalyzer().analyze(withLicenseFiles...)
	grc.logErrs(errs)
	return grc.renderer.render(os.Stdout, licPkgs...)
}

func (grc genericRenderCmd) envPerm(cmd *cobra.Command) [][]string {
	gooses, _ := cmd.Flags().GetStringSlice("gooses")
	var perms [][]string
	for _, goos := range gooses {
		perms = append(perms, []string{"GOOS=" + goos})
	}
	return perms
}

func (grc *genericRenderCmd) logErrs(errs []error) {
	for _, err := range errs {
		if err != nil {
			log.Println(err.Error())
		}
	}
}
