package cmd

import (
	"fmt"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var blue = color.New(color.Bold, color.FgBlue)

var listLogsCmd = &cobra.Command{
	Use:     "logs",
	Short:   "List recent FRITZ!BOX logs",
	Long:    "List the log statements/events from the FRITZ!Box. Logs may be subject to log rotation by the FRITZ!Box.",
	Example: "fritzctl list logs",
	RunE:    listLogs,
}

func init() {
	listCmd.AddCommand(listLogsCmd)
}

func listLogs(cmd *cobra.Command, args []string) error {
	c := clientLogin()
	f := fritz.Internal(c)
	logs, err := f.ListLogs()
	assert.NoError(err, "cannot obtain logs:", err)
	logger.Success("Obtained log messages:")
	printLogs(logs)
	return nil
}

func printLogs(logs *fritz.MessageLog) {
	for _, m := range logs.Messages {
		printLog(&m)
	}
}

func printLog(m *fritz.Message) {
	text := (*m)[0]
	if len(text) >= 17 {
		blue.Print(text[:17])
		fmt.Println(text[17:])
	} else {
		fmt.Println(text)
	}
}
