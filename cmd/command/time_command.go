package command

import (

	//	"github.com/google/uuid"
	"github.com/talbx/go-clockodo/cmd/command/timeprocessing"
	"github.com/talbx/go-clockodo/cmd/util"
)

type TimeCommand struct{}

func (cmd TimeCommand) Execute() {
	var processor timeprocessing.TimeProcessor = timeprocessing.CreateCommandFactory().CreateTimeProcessor()
	round, _ := util.GetFlags().GetBool("round")
	processor.Process(round)
}
