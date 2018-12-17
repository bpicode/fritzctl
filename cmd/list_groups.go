package cmd

import (
	"os"
	"sort"
	"strings"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/internal/console"
	"github.com/bpicode/fritzctl/internal/stringutils"
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

func listGroups(_ *cobra.Command, _ []string) error {
	list := mustList()
	printGroups(list)
	return nil
}

type byGroupName []fritz.DeviceGroup

// Len returns the length of the slice.
func (p byGroupName) Len() int {
	return len(p)
}

// Swap exchanges elements in the slice.
func (p byGroupName) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Less compares group names.
func (p byGroupName) Less(i, j int) bool {
	return p[i].Group.Name < p[j].Group.Name
}

func printGroups(list *fritz.Devicelist) {
	groups := list.DeviceGroups()
	sort.Sort(byGroupName(groups))
	table := console.NewTable(console.Headers("NAME", "MEMBERS", "MASTER", "PRESENT", "STATE", "TEMP (MEAS/WANT/SAV/COMF) [°C]"))
	for _, g := range groups {
		table.Append(groupColumns(g, list))
	}
	table.Print(os.Stdout)
}

func groupColumns(group fritz.DeviceGroup, list *fritz.Devicelist) []string {
	return []string{
		group.Group.Name,
		strings.Join(memberNames(group, list), ", "),
		masterName(group, list),
		console.IntToCheckmark(group.Group.Present),
		console.StringToCheckmark(group.Group.Switch.State),
		joinTemperatures(group.Group.Thermostat),
	}
}

func memberNames(group fritz.DeviceGroup, list *fritz.Devicelist) []string {
	var names []string
	for _, id := range group.Group.Members() {
		if d, ok := list.DeviceWithID(id); ok {
			names = append(names, d.Name)
		}
	}
	sort.Strings(names)
	return names
}

func masterName(group fritz.DeviceGroup, list *fritz.Devicelist) string {
	master, ok := list.DeviceWithID(group.Group.GroupInfo.MasterDeviceID)
	if !ok {
		return "(none)"
	}
	return master.Name
}

func joinTemperatures(th fritz.Thermostat) string {
	return strings.Join([]string{
		valueOrQm(th.FmtMeasuredTemperature),
		valueOrQm(th.FmtGoalTemperature),
		valueOrQm(th.FmtSavingTemperature),
		valueOrQm(th.FmtComfortTemperature)},
		"/")
}

func valueOrQm(f func() string) string {
	yellowQm := console.Yellow("?")
	return stringutils.DefaultIfEmpty(f(), yellowQm)
}
