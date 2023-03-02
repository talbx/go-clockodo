package command

import (
	"github.com/talbx/go-clockodo/cmd/intercept"
)

type Factory interface {
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

var instance Factory

func CreateCommandFactory() *Factory {
	if instance == nil {
		instance = ClockodoCommandFactory{}
	}
	return &instance
}
