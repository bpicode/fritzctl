package cmd

import (
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use: "switch [subcommand]",
}

func init() {
	RootCmd.AddCommand(switchCmd)
	switchOnCmd.Commands()
}
