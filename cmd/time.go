package cmd

import (
	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "shows your clocked time",
	Run: Process,
}

func init() {
	rootCmd.AddCommand(timeCmd)
}
