package concurrent

import (
	"fmt"
	. "github.com/talbx/go-clockodo/pkg/model"
	"log/slog"
	"strings"
	"sync"
	"time"
)

type ConcurrentDayEntriesByCustomerAggregator struct {
	TimeEntries *map[string][]TimeEntry
	CustomerIds []int
}

func (c ConcurrentDayEntriesByCustomerAggregator) Aggregate() *map[string][]DayByCustomer {
	var mappyNew = CreateDBC()
	for _, customerId := range c.CustomerIds {
		p_cdbc(c.TimeEntries, customerId, mappyNew)
	}
	return mappyNew.Get()
}

func p_cdbc(mappy *map[string][]TimeEntry, customerId int, syncMap *TimeEntrySyncMap[DayByCustomer]) {
	wg := sync.WaitGroup{}
	for w, v := range *mappy {
		wg.Add(1)
		go tryfunc(w, v, &wg, syncMap, customerId)
	}
	slog.Debug(fmt.Sprintf("[CreateDayByCustomer] using %v worker threads to rebuild models\n", len(*mappy)))
	wg.Wait()
	slog.Debug("[CreateDayByCustomer] done rebuilding models\n")
}

func tryfunc(w string, v []TimeEntry, wg *sync.WaitGroup, m *TimeEntrySyncMap[DayByCustomer], customerId int) {
	dbc := DayByCustomer{CustomerId: customerId}
	addUpTimeAndTasks(v, customerId, &dbc)
	if dbc.Tasks != "" {
		cleanTasks := RemoveDuplicateStr(strings.Split(dbc.Tasks, ","))
		dbc.AggregatedTasks = strings.Join(cleanTasks, ",")
		h, _ := time.ParseDuration(fmt.Sprintf("%ds", dbc.TotalTime))
		dbc.AggregatedTime = h.String()
		m.AppendNonExistent(w, dbc)
	}
	wg.Done()
}

func addUpTimeAndTasks(v []TimeEntry, customerId int, dbc *DayByCustomer) {
	for _, te := range v {
		if te.CustomerId == customerId {
			dbc.TotalTime += te.Duration
			dbc.Tasks += ", " + te.Description
		}
	}
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value && item != "" && item != " " {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
