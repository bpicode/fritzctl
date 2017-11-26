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

func printVersion(_ *cobra.Command, _ []string) error {
	fmt.Printf("%s (revision %s)\n", config.Version, config.Revision)
	return nil
}
