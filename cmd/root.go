package cmd

import (
	"github.com/bpicode/fritzctl/flags"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any sub-commands.
var RootCmd = &cobra.Command{
	Use: "fritzctl",
}

func init() {
	cobra.OnInitialize()
	RootCmd.PersistentFlags().Var(&flags.Loglevel{}, "loglevel", "logging verbosity")
}
