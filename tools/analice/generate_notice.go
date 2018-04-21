package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var noticeCmd = &cobra.Command{
	Use:  "notice /path/to/project",
	RunE: genNotice,
}

func init() {
	generateCmd.AddCommand(noticeCmd)
}

func genNotice(_ *cobra.Command, args []string) error {
	dir := projectDir(args)
	p, err := getProjector().project(dir)
	if err != nil {
		return fmt.Errorf("failed to get project dependencies: %v", err)
	}
	m := runEngine(p)
	n := noticeWriter{}
	n.writeTo(os.Stdout, m)
	return nil
}

func runEngine(p *project) map[string][]dependency {
	bowEngine := bowLicenseEngine{}
	a := newAnalyzer(&bowEngine)
	a.run(*p)
	lvd := bowEngine.licsVsDeps
	m := make(map[string][]dependency)
	for k, v := range lvd {
		m[k.name] = append(m[k.name], v...)
	}
	return m
}

func newAnalyzer(las ...lineAnalyzer) analyzer {
	a := analyzer{}
	for _, la := range las {
		a.register(la)
	}
	return a
}
