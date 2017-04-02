package cmd

import (
	"github.com/mitchellh/cli"
	"github.com/bpicode/fritzctl/completion"
)

type completionBashCommand struct {
	cli *cli.CLI
}

func (cmd *completionBashCommand) Help() string {
	return `Outputs fritzctl shell completion for the given shell (bash)
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
Additionally, you may want to output completion to a file and source in your .bashrc`
}

func (cmd *completionBashCommand) Synopsis() string {
	return "outputs fritzctl shell completion for the given shell (bash)"
}

func (cmd *completionBashCommand) Run(args []string) int {
	commands := make([]string, len(cmd.cli.Commands))
	for command := range cmd.cli.Commands {
		commands = append(commands, command)
	}
	bash := completion.BourneAgainShell()
	bash.Print()
	return 0
}

// CompletionBash returns a factory creating commands for the
// bash completion bindings of this app.
func CompletionBash(c *cli.CLI) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &completionBashCommand{cli: c}, nil
	}
}
