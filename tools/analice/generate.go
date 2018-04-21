package main

import "github.com/spf13/cobra"

var generateCmd = &cobra.Command{
	Use: "generate [subcommand]",
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

func projectDir(args []string) string {
	var dir = "."
	if len(args) > 0 {
		dir = args[0]
	}
	return dir
}

func getProjector() projector {
	return &depLockProjector{}
}
