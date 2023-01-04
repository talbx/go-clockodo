package stop

import "github.com/google/uuid"

type StopProcessor struct{}

func (processor StopProcessor) Process() (string, string){
	pid := uuid.New().String()

	return "", pid
}