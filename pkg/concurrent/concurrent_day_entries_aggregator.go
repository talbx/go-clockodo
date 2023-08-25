package concurrent

import (
	"fmt"
	. "github.com/talbx/go-clockodo/pkg/model"
	"log/slog"
	"sync"
	"time"
)

type ConcurrentDayEntriesAggregator struct {
	TimeEntries TimeEntriesResponse
}

func (c ConcurrentDayEntriesAggregator) Aggregate() *map[string][]TimeEntry {
	var weekDayMap = CreateTE()
	wg := sync.WaitGroup{}
	workerCount := 0
	start := time.Now()
	f := func(entry TimeEntry, m *TimeEntrySyncMap[TimeEntry], wg *sync.WaitGroup, i int) {
		defer wg.Done()
		startTime, _ := time.Parse("2006-01-02T15:04:05Z", entry.StartTime)
		m.AppendNonExistent(startTime.Format("02.01.2006"), entry)
	}
	for i, te := range c.TimeEntries.Entries {
		wg.Add(1)
		workerCount++
		go f(te, weekDayMap, &wg, i)
	}
	slog.Debug(fmt.Sprintf("[GroupEntriesByDay] using %v worker threads to match time entries its weekdays\n", workerCount))
	wg.Wait()
	end := time.Now()
	t := end.Sub(start)
	slog.Debug(fmt.Sprintf("[GroupEntriesByDay] done matching time entries its weekdays in %v microseconds (%v seconds)\n", t.Microseconds(), t.Seconds()))
	return weekDayMap.Get()
}
