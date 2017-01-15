package cmd

import (
	"fmt"

	"strings"

	"github.com/bpicode/fritzctl/assert"
	"github.com/bpicode/fritzctl/conv"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/mitchellh/cli"
)

type listInetstatsCommand struct {
}

func (cmd *listInetstatsCommand) Help() string {
	return "Get recent internet upstream/downstream statistics from the FRITZ!Box."
}

func (cmd *listInetstatsCommand) Synopsis() string {
	return "get recent internet statistics"
}

func (cmd *listInetstatsCommand) Run(args []string) int {
	c := clientLogin()
	f := fritz.New(c)
	stats, err := f.InternetStats()
	assert.NoError(err, "cannot obtain internet stats:", err)
	logger.Success("Obtained recent upstream/downstream time series:\n")
	printTrafficData(stats)
	return 0
}

// ListInetstats is a factory creating commands for commands listing internet statistics.
func ListInetstats() (cli.Command, error) {
	p := listInetstatsCommand{}
	return &p, nil
}

func printTrafficData(data *fritz.TrafficMonitoringData) {
	kbps := data.KiloBitsPerSecond()
	printSlice("Downstream/internet       [kb/s]: ", kbps.DownstreamInternet)
	printSlice("Downstream/media          [kb/s]: ", kbps.DownStreamMedia)
	printSlice("Upstream/low priority     [kb/s]: ", kbps.UpstreamLowPriority)
	printSlice("Upstream/default priority [kb/s]: ", kbps.UpstreamDefaultPriority)
	printSlice("Upstream/high priority    [kb/s]: ", kbps.UpstreamHighPriority)
	printSlice("Upstream/realtime         [kb/s]: ", kbps.UpstreamRealtime)
}

func printSlice(pre string, data []float64) {
	strs := conv.Float64Slice(data).String('f', 2)
	fmt.Println(pre, strings.Join(strs, " "))
}
