package cmd

import (
	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:     "sessionid",
	Short:   "Check if the FRITZ!Box responds",
	Long:    "Attempt to contact the FRITZ!Box by trying to solve the login challenge.",
	Example: "fritzctl ping",
	RunE:    ping,
}

func init() {
	RootCmd.AddCommand(pingCmd)
}

func ping(cmd *cobra.Command, args []string) error {
	clientLogin()
	logger.Success("Success! FRITZ!Box seems to be alive!")
	return nil
}
