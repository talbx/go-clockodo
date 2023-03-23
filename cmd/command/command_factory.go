package command

import (
	"github.com/talbx/go-clockodo/cmd/command/timeprocessing"
	"github.com/talbx/go-clockodo/pkg/intercept"
)

type Factory interface {
	Create(cmd string) timeprocessing.TimeProcessor
}

type Command interface {
	Execute()
}

type ClockodoCommandFactory struct{}

func (factory ClockodoCommandFactory) Create(cmd string) timeprocessing.TimeProcessor {
	intercept.ConfigReaderInterceptor{}.Intercept()
	return timeprocessing.WeekProcessor{}
}

var instance Factory

func CreateCommandFactory() *Factory {
	if instance == nil {
		instance = ClockodoCommandFactory{}
	}
	return &instance
}
