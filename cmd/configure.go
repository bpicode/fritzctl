package cmd

import (
	"io"
	"os"

	"github.com/bpicode/fritzctl/config"
	"github.com/spf13/cobra"
)

var (
	configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Configure fritzctl",
		Long: `Walk through the configuration of fritzctl interactively.
Configuration file is saved at the end of the survey.",
Run fritzctl with administrator privileges if the configuration file cannot be saved by a normal user.`,
		Example: "fritzctl configure",
		RunE:    configure,
	}
	configReaderSrc io.Reader = os.Stdin
)

func init() {
	RootCmd.AddCommand(configureCmd)
}

func configure(cmd *cobra.Command, args []string) error {
	configurer := config.NewConfigurer()
	configurer.Greet()
	cfg, err := configurer.Obtain(configReaderSrc)
	assertNoErr(err, "error obtaining configuration")
	err = cfg.Write()
	assertNoErr(err, "error writing configuration file")
	return nil
}
