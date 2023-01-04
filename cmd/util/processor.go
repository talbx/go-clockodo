package util

// Simple Processor handling any process
// always producing the cli output and a processId for reference.
type Processor interface {
	Process() (result string, processId string)
}