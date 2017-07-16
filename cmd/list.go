package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use: "list",
}

func init() {
	RootCmd.AddCommand(listCmd)
}
