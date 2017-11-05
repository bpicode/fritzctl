package cmd_test

import (
	"fmt"

	"github.com/bpicode/fritzctl/cmd"
)

// RootCmd is the only exported variable. Sub-commands are unexported children of RootCmd or other unexported
// sub-commands.
func Example() {
	use := cmd.RootCmd.Use
	fmt.Println(use)
	// output: fritzctl [subcommand]
}
