package cmd

import (
	"fmt"
	"os"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/internal/console"
	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
)

var listCallsCmd = &cobra.Command{
	Use:     "phonecalls",
	Short:   "List recent phone calls",
	Long:    "List recently made phone calls.",
	Example: "fritzctl list phonecalls",
	RunE:    listCalls,
}

var callTypeVsDescription = map[string]string{
	"1": "incoming",
	"2": "missed",
	"3": "rejected",
	"4": "outgoing",
}

func init() {
	listCmd.AddCommand(listCallsCmd)
}

func listCalls(_ *cobra.Command, _ []string) error {
	c := clientLogin()
	f := fritz.NewPhone(c)
	calls, err := f.Calls()
	assertNoErr(err, "cannot obtain phone calls")
	logger.Success("Recent phone calls:\n")
	printCalls(calls)
	return nil
}

func printCalls(calls []fritz.Call) {
	table := console.NewTable(console.Headers("TYPE", "DATE", "CALLER", "DURATION"))
	for _, c := range calls {
		table.Append(callColumns(c))
	}
	table.Print(os.Stdout)
}

func callColumns(c fritz.Call) []string {
	return []string{
		callTypeVsDescription[c.Type],
		c.Date,
		fmt.Sprintf("%s %s", c.Caller, c.PhoneNumber),
		c.Duration,
	}
}
