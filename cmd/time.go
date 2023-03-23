package cmd

import (
	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "shows your clocked time",
	Run:   Process,
}

func init() {
	timeCmd.Flags().IntP("last", "l", 0, "l1 = last week, l2 = 2 weeks ago. default l0 = now")
	rootCmd.AddCommand(timeCmd)
}
