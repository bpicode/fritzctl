package cmd

import (
	"github.com/bpicode/fritzctl/cmd/jsonapi"
	"github.com/bpicode/fritzctl/fritz"
	"github.com/spf13/cobra"
)

func selectFmt(cmd *cobra.Command, ds []fritz.Device, defaultF func([]fritz.Device) interface{}) interface{} {
	switch cmd.Flag("output").Value.String() {
	case "json":
		return jsonapi.NewMapper().Convert(ds)
	default:
		return defaultF(ds)
	}
}
