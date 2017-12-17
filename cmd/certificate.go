package cmd

import (
	"github.com/spf13/cobra"
)

var certCmd = &cobra.Command{
	Use:   "certificate",
	Short: "See subcommands",
	Long:  "See subcommands. Run with --help to list the available commands.",
}

func init() {
	RootCmd.AddCommand(certCmd)
}
