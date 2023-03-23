package concurrent

import "time"

type ConcurrentAggegator interface {
	Aggregate() *map[time.Weekday][]interface{}
}
