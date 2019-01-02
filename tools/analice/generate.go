package main

import (
	"runtime"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use: "generate [subcommand]",
}

func init() {
	generateCmd.PersistentFlags().Bool("tests", true, "include test dependencies in license analysis")
	generateCmd.PersistentFlags().StringSlice("gooses", []string{runtime.GOOS}, "run analysis for these GOOS values")
	rootCmd.AddCommand(generateCmd)
}
