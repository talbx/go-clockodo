package cmd

import (
	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "shows your clocked time for the current week",
	Run:   Process,
}

func init() {
	rootCmd.AddCommand(todayCmd)
}
