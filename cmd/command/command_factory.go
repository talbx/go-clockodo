package command

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/talbx/go-clockodo/cmd/command/timeprocessing"
	"github.com/talbx/go-clockodo/pkg/intercept"
	"github.com/talbx/go-clockodo/pkg/util"
)

type Factory interface {
	Create(cmd string) util.Processor
}

type Command interface {
	Execute()
}

type ClockodoCommandFactory struct{}

func (factory ClockodoCommandFactory) Create(cmd string) util.Processor {
	err := intercept.ConfigReaderInterceptor{}.Intercept()
	if err != nil {
		slog.Error("no config.yaml could be found. please provide one")
		os.Exit(1) // Handle errors reading the config file
	}
	we, _ := util.GetFlags().GetBool("withEarnings")
	slog.Debug(fmt.Sprintf("with earnings %v", we))
	intercept.ClockodoConfig.WithRevenue = we
	slog.Debug(fmt.Sprintf("%+v", intercept.ClockodoConfig))

	if cmd == "today" {
		return timeprocessing.TodayProcessor{}
	}

	return timeprocessing.TimeProcessor{}
}

var instance Factory

func CreateCommandFactory() *Factory {
	if instance == nil {
		instance = ClockodoCommandFactory{}
	}
	return &instance
}
