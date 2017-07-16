package cmd

import (
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use: "switch",
}

func init() {
	RootCmd.AddCommand(switchCmd)
}
