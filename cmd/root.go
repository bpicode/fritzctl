package cmd

import (
	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any sub-commands.
var RootCmd = &cobra.Command{
	Use:   "fritzctl [subcommand]",
	Short: "A lightweight, easy to use console client for the AVM FRITZ!Box Home Automation",
	Long: "fritzctl is a command line client for the AVM FRITZ!Box primarily focused on the AVM Home Automation HTTP Interface. " +
		"For recent developments and releases visit https://github.com/bpicode/fritzctl. " +
		"For the vendor description visit https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf.",
}

func init() {
	cobra.OnInitialize()
	RootCmd.PersistentFlags().Var(&logger.Level{}, "loglevel", "logging verbosity")
	RootCmd.InitDefaultHelpFlag()
	RootCmd.InitDefaultHelpCmd()
}
