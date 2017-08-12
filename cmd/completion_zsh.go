package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Outputs fritzctl shell completion for the given shell (zsh)",
	Long: `Outputs fritzctl shell completion for the given shell (zsh) to stdout.
Completion functions for need to be stored a file called _fritzctl and this file needs to be placed in a directory listed in the $fpath variable.
Afterwards restart the shell or run

autoload -U compinit && compinit

to make zsh aware of the changes.`,
	Example: "sudo sh -c '/path/to/fritzctl completion zsh > /usr/share/zsh/vendor-completions/_fritzctl' && autoload -U compinit && compinit",
	RunE:    completionZsh,
}

func init() {
	completionCmd.AddCommand(completionZshCmd)
}

func completionZsh(cmd *cobra.Command, args []string) error {
	return RootCmd.GenZshCompletion(os.Stdout)
}
