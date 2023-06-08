package cmd

import (
	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:   "week",
	Short: "shows your clocked time for the current week",
	Run:   Process,
}

func init() {
	timeCmd.Flags().IntP("last", "l", 0, "l1 = last week, l2 = 2 weeks ago. default l0 = now")
	timeCmd.Flags().BoolP("withEarnings", "e", false, "if true, earnings are calculated based on config")
	rootCmd.AddCommand(timeCmd)
}
