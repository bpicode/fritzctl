package cmd

import (
	"github.com/spf13/cobra"
)

var manifestCmd = &cobra.Command{
	Use: "manifest",
}

func init() {
	RootCmd.AddCommand(manifestCmd)
}
