package timeprocessing

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/talbx/go-clockodo/cmd/util"
	"golang.org/x/exp/slices"
)

type WeekProcessor struct{}

type WeekByCustomer struct {
	Customer        string
	Tasks           string
	TotalTime       int
	AggregatedTime  string
	CustomerId      int
	AggregatedTasks string
}

type DayByCustomer struct {
	Customer        string
	Tasks           string
	TotalTime       int
	AggregatedTime  string
	CustomerId      int
	AggregatedTasks string
}

func (p WeekProcessor) Process(round bool) {
	var entriesRoot string = "v2/entries"

	now := time.Now()
	monday := getMonday(now)
	timeUntil := EOB(now).Format("2006-01-02T15:04:05Z")
	var repo util.TimeEntriesResponse
	query := fmt.Sprintf("%s?time_since=%s&time_until=%s", entriesRoot, monday, timeUntil)
	util.CallApi(query, &repo)

	mappy := groupEntriesByDay(repo)

	custIds := make([]int, 0)
	for _, e := range repo.Entries {
		if !slices.Contains(custIds, e.CustomerId) {
			custIds = append(custIds, e.CustomerId)
		}
	}

	mappyNew := groupDayEntriesByCustomer(mappy, custIds)
	cache := util.CustomerNameCache{}

	for k := range mappyNew {
		channel := make(chan DayByCustomer)
		for _, customerDay := range mappyNew[k] {
			go getAndAddCustomerNamesToEntries(customerDay, cache, channel)
		}
		resultList := make([]DayByCustomer, 0)
		for i := 0; i < len(mappyNew[k]); i++ {
			resultList = append(resultList, <-channel)
		}
		mappyNew[k] = resultList
	}

	var currentClock string = "v2/clock"
	var clo util.ClockResponse
	util.CallApi(currentClock, &clo)

	//fmt.Printf("%v+", c2s)
	Output(mappyNew, clo)
}

func getAndAddCustomerNamesToEntries(v DayByCustomer, cache util.CustomerNameCache, c1 chan DayByCustomer) {
	v.Customer = cache.Get(v.CustomerId)
	c1 <- v
}

func groupDayEntriesByCustomer(mappy map[time.Weekday][]util.TimeEntry, custIds []int) map[time.Weekday][]DayByCustomer {
	mappyNew := make(map[time.Weekday][]DayByCustomer)
	for _, customerId := range custIds {
		dbc := DayByCustomer{CustomerId: customerId}
		for w, v := range mappy {
			for _, te := range v {
				if te.CustomerId == customerId {
					dbc.TotalTime += te.Duration
					dbc.Tasks += ", " + te.Description
				}
			}
			cleanTasks := RemoveDuplicateStr(strings.Split(dbc.Tasks, ","))
			dbc.AggregatedTasks = strings.Join(cleanTasks, ",")
			h, _ := time.ParseDuration(fmt.Sprintf("%ds", dbc.TotalTime))
			dbc.AggregatedTime = h.String()
			mappyNew[w] = append(mappyNew[w], dbc)
		}
	}
	return mappyNew
}

func groupEntriesByDay(repo util.TimeEntriesResponse) map[time.Weekday][]util.TimeEntry {
	weekdayMap := make(map[time.Weekday][]util.TimeEntry)
	weekdays := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}
	for _, weekday := range weekdays {
		for _, te := range repo.Entries {
			startTime, _ := time.Parse("2006-01-02T15:04:05Z", te.StartTime)
			if weekday == startTime.Weekday() {
				weekdayMap[weekday] = append(weekdayMap[weekday], te)
			}
		}
	}
	return weekdayMap
}

func Output(mappy map[time.Weekday][]DayByCustomer, clock util.ClockResponse) {

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Customer", "Tasks", "Times"})
    rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	taskCount := 0
	tt := 0
	for key := range mappy {
		for _, entry := range mappy[key] {
			taskCount += len(strings.Split(entry.AggregatedTasks, ","))
			tt += entry.TotalTime
			t.AppendRow(table.Row{key, entry.CustomerId, entry.Customer, entry.AggregatedTasks, entry.AggregatedTime}, rowConfigAutoMerge)
		}
	}

	t.SetColumnConfigs([]table.ColumnConfig{
        {Number: 1, AutoMerge: true},
        {Number: 2, AutoMerge: true},
        {Number: 3, AutoMerge: true},
        {Number: 4, AutoMerge: true},
        {Number: 5, AutoMerge: true},
        {Number: 6, AutoMerge: true},
    })
	th, tm := DurationToHM(tt)
	t.AppendSeparator()
	t.SetStyle(table.StyleLight)
	t.AppendFooter(table.Row{"TOTAL", "", "", fmt.Sprintf("Total tasks: %v", taskCount), fmt.Sprintf("%v:%v", AddLeadingZero(th), AddLeadingZero(tm))})
	t.Render()

	h, m := ProcessClock(&clock)
	log.Printf("Also, you have a task running for %vh:%vm right now.\n", h, m)
}

func getMonday(t time.Time) string {
	if t.Weekday() == time.Monday {
		return SOB(t).Format("2006-01-02T15:04:05Z")
	}
	return getMonday(t.AddDate(0, 0, -1))
}
