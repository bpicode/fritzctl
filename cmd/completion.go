package cmd

import (
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use: "completion",
}

func init() {
	RootCmd.AddCommand(completionCmd)
}
