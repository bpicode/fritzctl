package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/logger"
	"github.com/spf13/cobra"
)

var listInetstatsCmd = &cobra.Command{
	Use:     "inetstats",
	Short:   "Get recent internet statistics",
	Long:    "Get recent internet upstream/downstream statistics from the FRITZ!Box.",
	Example: "fritzctl list inetstats",
	RunE:    listInetstats,
}

func init() {
	listCmd.AddCommand(listInetstatsCmd)
}

func listInetstats(_ *cobra.Command, _ []string) error {
	c := clientLogin()
	f := fritz.NewInternal(c)
	stats, err := f.InternetStats()
	assertNoErr(err, "cannot obtain internet stats")
	logger.Success("Recent upstream/downstream time series:\n")
	printTrafficData(stats)
	return nil
}

func printTrafficData(data *fritz.TrafficMonitoringData) {
	kbps := data.KiloBitsPerSecond()
	printSlice("Downstream/internet       [kb/s]: ", kbps.DownstreamInternet)
	printSlice("Downstream/media          [kb/s]: ", kbps.DownStreamMedia)
	printSlice("Downstream/guest          [kb/s]: ", kbps.DownStreamGuest)
	printSlice("Upstream/low priority     [kb/s]: ", kbps.UpstreamLowPriority)
	printSlice("Upstream/default priority [kb/s]: ", kbps.UpstreamDefaultPriority)
	printSlice("Upstream/high priority    [kb/s]: ", kbps.UpstreamHighPriority)
	printSlice("Upstream/realtime         [kb/s]: ", kbps.UpstreamRealtime)
	printSlice("Upstream/guest            [kb/s]: ", kbps.UpstreamGuest)
}

func printSlice(pre string, data []float64) {
	strs := float64Slice(data).formatFloats('f', 2)
	fmt.Println(pre, strings.Join(strs, " "))
}

type float64Slice []float64

func (fs float64Slice) formatFloats(format byte, prec int) []string {
	var strs []string
	for _, f := range fs {
		strs = append(strs, strconv.FormatFloat(f, format, prec, 64))
	}
	return strs
}
