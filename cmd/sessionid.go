package cmd

import (
	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
)

var sessionIDCmd = &cobra.Command{
	Use:   "sessionid",
	Short: "Obtain a session ID",
	Long: `Obtain a session ID by solving the FRITZ!Box login challenge. The session ID can be used for subsequent requests until it gets invalidated.
Visit https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AVM_Technical_Note_-_Session_ID.pdf for more information.`,
	Example: "fritzctl sessionid",
	RunE:    sessionID,
}

func init() {
	RootCmd.AddCommand(sessionIDCmd)
}

func sessionID(cmd *cobra.Command, args []string) error {
	client := clientLogin()
	logger.Success("Successfully obtained session ID:", client.SessionInfo.SID)
	return nil
}
