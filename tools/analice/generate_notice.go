package main

import (
	"github.com/spf13/cobra"
)

func init() {
	generateCmd.AddCommand(noticeCmd)
}

var noticeCmd = &cobra.Command{
	Use: "notice /path/to/project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return genericRenderCmd{renderer: newNoticeFileRenderer(), scanner: newDepScanner()}.run(cmd, args)
	},
}
