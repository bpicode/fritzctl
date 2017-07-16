package cmd

import (
	"github.com/spf13/cobra"
)

var manifestCmd = &cobra.Command{
	Use:   "manifest [subcommand]",
	Short: "See subcommands",
	Long:  "See subcommands. Run with --help to list the available commands.",
}

func init() {
	RootCmd.AddCommand(manifestCmd)
}
