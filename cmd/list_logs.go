package cmd

import (
	"fmt"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/fatih/color"
	"github.com/mitchellh/cli"
)

var (
	blue = color.New(color.Bold, color.FgBlue)
)

type listLogsCommand struct {
}

func (cmd *listLogsCommand) Help() string {
	return "List the log statements/events from the FRITZ!Box. Logs may be subject to log rotation by the FRITZ!Box."
}

func (cmd *listLogsCommand) Synopsis() string {
	return "list recent FRITZ!BOX logs"
}

func (cmd *listLogsCommand) Run(args []string) int {
	c := clientLogin()
	f := fritz.New(c)
	logs, err := f.ListLogs()
	assert.NoError(err, "cannot obtain logs:", err)
	logger.Success("Obtained log messages:\n")
	printLogs(logs)
	return 0
}

func printLogs(logs *fritz.MessageLog) {
	for _, m := range logs.Messages {
		printLog(&m)
	}
}

func printLog(m *fritz.Message) {
	text := m.Text
	if len(text) >= 17 {
		blue.Print(text[:17])
		fmt.Println(text[17:])
	} else {
		fmt.Println(text)
	}
}

// ListLogs is a factory creating commands for commands listing logs.
func ListLogs() (cli.Command, error) {
	p := listLogsCommand{}
	return &p, nil
}
