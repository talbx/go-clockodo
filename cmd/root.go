package cmd

import (
	"github.com/talbx/go-clockodo/pkg/util"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/talbx/go-clockodo/cmd/command"
)

var rootCmd = &cobra.Command{
	Use:   "go-clockodo",
	Run:   Process,
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func Process(cmd *cobra.Command, args []string) {
	s := time.Now()
	util.StoreFlags(cmd.Flags())
	var factory = *command.CreateCommandFactory()
	util.SugaredLogger.Infof("[CMD] processing command %v", cmd.Use)
	last, err := cmd.Flags().GetInt("last")
	if err != nil {
		util.SugaredLogger.Fatal(err)
	}
	err = factory.Create(cmd.Use).Process(last)
	if err != nil {
		util.SugaredLogger.Fatal(err)
	}
	e := time.Now()
	util.SugaredLogger.Infof("[CMD] the process is done and took %vms", e.Sub(s).Milliseconds())
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
