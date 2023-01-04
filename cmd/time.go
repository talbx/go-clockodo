package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	ttime "github.com/talbx/go-clockodo/cmd/time"
	"github.com/talbx/go-clockodo/cmd/util"
)

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "shows your clocked time today",
	Run: run,
}

func run(cmd *cobra.Command, args []string){
	var proc util.Processor = ttime.GetStatusProcessorImpl{}
	result, _ := proc.Process()
	fmt.Print(result)
}

func init() {
	rootCmd.AddCommand(timeCmd)
}
