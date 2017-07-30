package cmd

import (
	"github.com/bpicode/fritzctl/config"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure fritzctl",
	Long: `Walk through the configuration of fritzctl interactively.
Configuration file is saved at the end of the survey.",
Run fritzctl with administrator privileges if the configuration file cannot be saved by a normal user.`,
	Example: "fritzctl configure",
	RunE:    configure,
}

func init() {
	RootCmd.AddCommand(configureCmd)
}

func configure(cmd *cobra.Command, args []string) error {
	configurer := config.NewConfigurer()
	configurer.ApplyDefaults(config.Defaults())
	configurer.Greet()
	configurer.Obtain()
	err := configurer.Write()
	assertNoError(err, "error writing configuration file:", err)
	return nil
}
