package cmd

import (
	"github.com/spf13/cobra"
)

var monthCmd = &cobra.Command{
	Use:   "month",
	Short: "shows your clocked time for the current month",
	Run:   Process,
}

func init() {
	monthCmd.Flags().BoolP("withEarnings", "e", false, "if true, earnings are calculated based on config")
	rootCmd.AddCommand(monthCmd)
}
