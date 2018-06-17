package cmd

import (
	"os"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/internal/console"
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

func sessionID(_ *cobra.Command, _ []string) error {
	client := clientLogin()
	logger.Success("Successfully obtained session ID:", client.SessionInfo.SID)
	printGrants(client.SessionInfo.Rights)
	return nil
}

func printGrants(rights fritz.Rights) {
	table := console.NewTable(console.Headers("RIGHT", "R", "W"))
	for i, n := range rights.Names {
		table.Append(grantColumns(n, rights.AccessLevels[i]))
	}
	table.Print(os.Stdout)
}

func grantColumns(name, access string) []string {
	mayRead := console.Btoc(access == "1" || access == "2")
	mayWrite := console.Btoc(access == "2")
	return []string{name, mayRead.String(), mayWrite.String()}
}
