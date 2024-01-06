package main

import (
	"os"

	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use: "easybackup",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func main() {
	if err := cmdRoot.Execute(); err != nil {
		os.Exit(1)
	}
}
