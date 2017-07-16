package cmd

import (
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch [subcommand]",
	Short: "See subcommands",
	Long:  "See subcommands. Run with --help to list the available commands.",
}

func init() {
	RootCmd.AddCommand(switchCmd)
	switchOnCmd.Commands()
}
