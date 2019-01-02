package main

import (
	"github.com/spf13/cobra"
)

var copyrightCmd = &cobra.Command{
	Use: "copyright /path/to/project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return genericRenderCmd{renderer: newDebianCopyrightRenderer(), scanner: newDepScanner()}.run(cmd, args)
	},
}

func init() {
	generateCmd.AddCommand(copyrightCmd)
}
