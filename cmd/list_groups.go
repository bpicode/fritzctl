package cmd

import (
	"os"
	"strings"

	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/spf13/cobra"
)

var listGroupsCmd = &cobra.Command{
	Use:     "groups",
	Short:   "List the device groups",
	Long:    "List the device groups configured at the FRITZ!Box.",
	Example: "fritzctl list groups",
	RunE:    listGroups,
}

func init() {
	listCmd.AddCommand(listGroupsCmd)
}

func listGroups(cmd *cobra.Command, args []string) error {
	c := homeAutoClient()
	list, err := c.List()
	assertNoError(err, "cannot obtain data for smart home groups:", err)
	printGroups(list)
	return nil
}

func printGroups(list *fritz.Devicelist) {
	groups := list.DeviceGroups()
	table := console.NewTable(console.Headers("NAME", "MEMBERS", "PRESENT"))
	for _, g := range groups {
		table.Append(groupColumns(g, list))
	}
	table.Print(os.Stdout)
}

func groupColumns(group fritz.DeviceGroup, list *fritz.Devicelist) []string {
	return []string{
		group.Group.Name,
		strings.Join(memberNames(group, list), ", "),
		console.IntToCheckmark(group.Group.Present),
	}
}

func memberNames(group fritz.DeviceGroup, list *fritz.Devicelist) []string {
	var names []string
	for _, id := range group.Group.Members() {
		if d, ok := list.DeviceWithID(id); ok {
			names = append(names, d.Name)
		}
	}
	return names
}
