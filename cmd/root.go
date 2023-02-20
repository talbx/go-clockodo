package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/talbx/go-clockodo/cmd/command"
	"github.com/talbx/go-clockodo/cmd/util"
)

var rootCmd = &cobra.Command{
	Use:   "go-clockodo",
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
	util.StoreFlags(cmd.Flags())
	var factory = *command.CreateCommandFactory()
	factory.Create(cmd.Use).Execute()
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
