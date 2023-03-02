package timeprocessing

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/talbx/go-clockodo/cmd/util"
	"golang.org/x/exp/slices"
)

type WeekProcessor struct{}

func (p WeekProcessor) Process(bool) {
	var entriesRoot = "v2/entries"

	now := time.Now()
	monday := getMonday(now)
	timeUntil := EOB(now).Format("2006-01-02T15:04:05Z")
	var repo util.TimeEntriesResponse
	query := fmt.Sprintf("%s?time_since=%s&time_until=%s", entriesRoot, monday, timeUntil)
	util.CallApi(query, &repo)

	mappy := util.GroupEntriesByDay(repo)
	custIds := extractCustomerIdsFromTimeEntries(repo)
	mappyNew := util.GroupDayEntriesByCustomer(mappy, custIds)
	veryNewMap := enhanceMapWithCustomerNames(mappyNew)
	clo := getCurrentClock()

	//fmt.Printf("%v+", c2s)
	Output(veryNewMap, clo)
}

func getCurrentClock() util.ClockResponse {
	var currentClock = "v2/clock"
	var clo util.ClockResponse
	util.CallApi(currentClock, &clo)
	return clo
}

func extractCustomerIdsFromTimeEntries(repo util.TimeEntriesResponse) []int {
	custIds := make([]int, 0)
	for _, e := range repo.Entries {
		if !slices.Contains(custIds, e.CustomerId) {
			custIds = append(custIds, e.CustomerId)
		}
	}
	return custIds
}

func enhanceMapWithCustomerNames(mappyNew map[time.Weekday][]util.DayByCustomer) map[time.Weekday][]util.DayByCustomer {
	cache := util.CustomerNameCache{}

	veryNewMap := make(map[time.Weekday][]util.DayByCustomer)

	// for every weekday in that map
	for k := range mappyNew {
		channel := make(chan util.DayByCustomer)

		// go over every task entry
		for _, customerDay := range mappyNew[k] {
			// and fetch its customerName, publish it to a weekday-group channel
			go getAndAddCustomerNamesToEntries(customerDay, cache, channel)
		}

		// collect all updated task entries for each week day
		resultList := make([]util.DayByCustomer, 0)
		for i := 0; i < len(mappyNew[k]); i++ {
			val := <-channel
			resultList = append(resultList, val)
		}
		veryNewMap[k] = resultList
	}
	return veryNewMap
}

func getAndAddCustomerNamesToEntries(v util.DayByCustomer, cache util.Cache, c1 chan util.DayByCustomer) {
	v.Customer = cache.Get(v.CustomerId)
	c1 <- v
}

func sortWeekdayMap(mappy map[time.Weekday][]util.DayByCustomer) map[time.Weekday][]util.DayByCustomer {
	nm := make(map[time.Weekday][]util.DayByCustomer)

	nm[time.Monday] = mappy[time.Monday]
	nm[time.Tuesday] = mappy[time.Tuesday]
	nm[time.Wednesday] = mappy[time.Wednesday]
	nm[time.Thursday] = mappy[time.Thursday]
	nm[time.Friday] = mappy[time.Friday]
	nm[time.Saturday] = mappy[time.Saturday]
	nm[time.Sunday] = mappy[time.Sunday]
	return nm
}

func Output(mappy map[time.Weekday][]util.DayByCustomer, clock util.ClockResponse) {

	for _, v := range mappy {
		sort.Slice(v[:], func(i, j int) bool {
			return v[i].Customer < v[j].Customer
		})
	}

	sortedMappy := sortWeekdayMap(mappy)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Customer", "Tasks", "Times"})
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	taskCount := 0
	tt := 0

	for key := range sortedMappy {
		for _, entry := range sortedMappy[key] {
			taskCount += len(strings.Split(entry.AggregatedTasks, ","))
			tt += entry.TotalTime
			t.AppendRow(table.Row{key, entry.CustomerId, entry.Customer, entry.AggregatedTasks, entry.AggregatedTime}, rowConfigAutoMerge)
			t.AppendSeparator()
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
