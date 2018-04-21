package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "analice [subcommand]",
}

func init() {
	cobra.OnInitialize()
	rootCmd.InitDefaultHelpFlag()
	rootCmd.InitDefaultHelpCmd()
}

var exitOnErr = log.Fatalln

func main() {
	err := rootCmd.Execute()
	if err != nil {
		exitOnErr(err)
	}
}
