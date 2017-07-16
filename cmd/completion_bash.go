package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Outputs fritzctl shell completion for the given shell (bash)",
	Long: `Outputs fritzctl shell completion for the given shell (bash)
This depends on the bash-completion binary. Example installation instructions:
OS X:
	$ brew install bash-completion
	$ source $(brew --prefix)/etc/bash_completion
	$ fritzctl completion bash > ~/.fritzctl-completion
	$ source ~/.fritzctl-completion
Ubuntu:
	$ apt-get install bash-completion
	$ source /etc/bash-completion
	$ source <(fritzctl completion bash)
Additionally, you may want to output completion to a file and source in your .bashrc`,
	Example: "fritzctl completion bash",
	RunE:    completionBash,
}

func init() {
	completionCmd.AddCommand(completionBashCmd)
}

func completionBash(cmd *cobra.Command, args []string) error {
	RootCmd.GenBashCompletion(os.Stdout)
	return nil
}
