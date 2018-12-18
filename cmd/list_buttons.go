package cmd

import (
	"os"
	"time"

	"github.com/bpicode/fritzctl/cmd/printer"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/internal/console"
	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
)

var listButtonsCmd = &cobra.Command{
	Use:     "buttons",
	Short:   "List the smart-home buttons",
	Long:    "List the all smart-home devices recognized as pressable buttons.",
	Example: "fritzctl list buttons",
	RunE:    listButtons,
}

func init() {
	listButtonsCmd.Flags().StringP("output", "o", "", "specify output format")
	listCmd.AddCommand(listButtonsCmd)
}

func listButtons(cmd *cobra.Command, _ []string) error {
	devs := mustList()
	data := selectFmt(cmd, devs.Buttons(), buttonsTable)
	logger.Success("Device data:")
	printer.Print(data, os.Stdout)
	return nil
}

func buttonsTable(devs []fritz.Device) interface{} {
	table := console.NewTable(console.Headers(
		"NAME",
		"LAST PRESSED",
	))
	referenceTime := time.Now()
	for _, dev := range devs {
		columns := []string{dev.Name, dev.Button.FmtLastPressedCompact(referenceTime)}
		table.Append(columns)
	}
	return table
}
