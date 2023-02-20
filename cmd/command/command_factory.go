package command

import (
	"github.com/talbx/go-clockodo/cmd/intercept"
)

type CommandFactory interface {
	Create(cmd string) Command
}

type Command interface {
	Execute()
}

type ClockodoCommandFactory struct{}

func (factory ClockodoCommandFactory) Create(cmd string) Command {
	intercept.ConfigReaderInterceptor{}.Intercept()
	return TimeCommand{}
}

var instance CommandFactory

func CreateCommandFactory() *CommandFactory {
	if instance == nil {
		instance = ClockodoCommandFactory{}
	}
	return &instance
}
