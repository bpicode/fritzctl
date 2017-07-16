package cmd

import (
	"fmt"

	"github.com/bpicode/fritzctl/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of this application",
	Long:  "Print the version of this application as one-line semantic version string.",
	RunE:  printVersion,
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func printVersion(cmd *cobra.Command, args []string) error {
	fmt.Println(config.Version)
	return nil
}
