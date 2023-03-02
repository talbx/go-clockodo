package command

import (

	//	"github.com/google/uuid"
	"github.com/talbx/go-clockodo/cmd/command/timeprocessing"
	"github.com/talbx/go-clockodo/cmd/util"
)

type TimeCommand struct{}

func (cmd TimeCommand) Execute() {
	round, _ := util.GetFlags().GetBool("round")
	timeprocessing.WeekProcessor{}.Process(round)
}
