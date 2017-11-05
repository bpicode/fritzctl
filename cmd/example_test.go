package cmd_test

import (
	"github.com/bpicode/fritzctl/cmd"
)

// RootCmd is the only exported variable. Sub-commands are unexported children of RootCmd or other unexported
// sub-commands.
func Example() {
	cmd.RootCmd.Execute()
	// output:
}
