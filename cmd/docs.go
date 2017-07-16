package cmd

import "github.com/spf13/cobra"

var docCmd = &cobra.Command{
	Use:   "doc [subcommand]",
	Short: "See subcommands",
	Long:  "See subcommands. Run with --help to list the available commands.",
}

func init() {
	RootCmd.AddCommand(docCmd)
}
