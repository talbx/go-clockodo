package timeprocessing

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
	. "github.com/talbx/go-clockodo/cmd/command/cashprocessing"
	"github.com/talbx/go-clockodo/pkg/concurrent"
	"github.com/talbx/go-clockodo/pkg/intercept"
	. "github.com/talbx/go-clockodo/pkg/model"
	. "github.com/talbx/go-clockodo/pkg/util"

	"github.com/jedib0t/go-pretty/v6/table"

	"golang.org/x/exp/slices"
)

type WeekProcessor struct{}

func (p WeekProcessor) Process(mode string, last int) error {

	SugaredLogger.Infof("last %v\n", last)
	monday, sunday := findStartAndFinish(mode, last)
	var entriesRoot = "v2/entries"

	var repo TimeEntriesResponse
	query := fmt.Sprintf("%s?time_since=%s&time_until=%s", entriesRoot, monday, sunday)
	_, err := CallApi(query, &repo)

	if err != nil {
		SugaredLogger.Fatal(err)
	}

	custIds := extractCustomerIdsFromTimeEntries(repo)
	dayEntriesAggregator := concurrent.ConcurrentDayEntriesAggregator{TimeEntries: repo}
	dayEntries := dayEntriesAggregator.Aggregate()

	dayEntriesByCustomerAggregator := concurrent.ConcurrentDayEntriesByCustomerAggregator{TimeEntries: dayEntries, CustomerIds: custIds}
	groupedDayEntries := dayEntriesByCustomerAggregator.Aggregate()

	enhancer := concurrent.ConcurrentCustomerNameEnhancer{TimeEntries: *groupedDayEntries}
	result := enhancer.Aggregate()
	clo, err := getCurrentClock()

	if err != nil {
		SugaredLogger.Fatal(err)
	}

	//fmt.Printf("%v+", c2s)
	Output(*result, clo)
	return nil
}

func findStartAndFinish(mode string, last int) (string, string) {
	if mode == "month" {
		now := time.Now()
		y := now.Year()
		m := now.Month()
		start := time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
		return SOB(start).Format("2006-01-02T15:04:05Z"), EOB(now).Format("2006-01-02T15:04:05Z")
	}

	if last == 0 {
		return getMonday(time.Now()), getSunday(time.Now())
	}
	now := time.Now().AddDate(0, 0, -(7 * last))

	fmt.Println(now.String())
	monday := getMonday(now)
	sunday := getSunday(now)
	return monday, sunday
}

func getCurrentClock() (ClockResponse, error) {
	var currentClock = "v2/clock"
	var clo ClockResponse
	_, err := CallApi(currentClock, &clo)
	return clo, err
}

func extractCustomerIdsFromTimeEntries(repo TimeEntriesResponse) []int {
	custIds := make([]int, 0)
	for _, e := range repo.Entries {
		if !slices.Contains(custIds, e.CustomerId) {
			custIds = append(custIds, e.CustomerId)
		}
	}
	return custIds
}

func sortWeekdayMap(mappy map[string][]DayByCustomer) map[string][]DayByCustomer {
	r := make(map[string][]DayByCustomer, 0)
	keys := make([]string, 0, len(mappy))
	for k := range mappy {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		r[k] = mappy[k]
	}

	return r
}

func alterTime(entry *DayByCustomer) {
	var hs, m = 0, 0
	if strings.Contains(entry.AggregatedTime, "h") {
		hRest := strings.Split(entry.AggregatedTime, "h")
		mRest := strings.Split(hRest[1], "m")
		hs, _ = strconv.Atoi(hRest[0])
		m, _ = strconv.Atoi(mRest[0])

	} else if strings.Contains(entry.AggregatedTime, "m") {
		mRest := strings.Split(entry.AggregatedTime, "m")
		m, _ = strconv.Atoi(mRest[0])
	}
	r1, r2 := Round(hs, m)
	entry.RoundedTime = fmt.Sprintf("%v:%v", r1, r2)
}

func Output(mappy map[string][]DayByCustomer, clock ClockResponse) {
	for _, v := range mappy {
		sort.Slice(v[:], func(i, j int) bool {
			return v[i].Customer < v[j].Customer
		})
	}

	sortedMappy := sortWeekdayMap(mappy)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "ID", "Customer", "Tasks", "Times", "Revenue"})
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	taskCount := 0
	tt := 0
	totalRevenue := money.New(0, money.EUR)
	keys := make([]string, 0, len(sortedMappy))
	for k := range sortedMappy {
		keys = append(keys, k)
	}

	for _, key := range keys {
		for _, entry := range sortedMappy[key] {
			alterTime(&entry)
			entry.AggregatedRevenue = money.New(0, money.EUR)
			if intercept.ClockodoConfig.WithRevenue {
				CashProcess(&entry)
				totalRevenue, _ = totalRevenue.Add(entry.AggregatedRevenue)
			}
			taskCount += len(strings.Split(entry.AggregatedTasks, ","))
			tt += entry.TotalTime
			t.AppendRow(table.Row{key, entry.CustomerId, entry.Customer, entry.AggregatedTasks, fmt.Sprintf("(%v) - %v", entry.RoundedTime, entry.AggregatedTime), entry.AggregatedRevenue.Display()}, rowConfigAutoMerge)
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
		{Number: 7, AutoMerge: true},
	})

	if intercept.ClockodoConfig.WithRevenue && intercept.ClockodoConfig.Revenue.RevenueStyle == "AN" {
		RevenueToANRevenue(totalRevenue)
	}

	th, tm := DurationToHM(tt)
	t.AppendSeparator()
	t.SetStyle(table.StyleLight)
	t.AppendFooter(table.Row{"TOTAL", "", "", fmt.Sprintf("Total tasks: %v", taskCount), fmt.Sprintf("%v:%v", AddLeadingZero(th), AddLeadingZero(tm)), totalRevenue.Display()})
	t.Render()

	h, m := ProcessClock(&clock)
	SugaredLogger.Infof("Also, you have a task running for %vh:%vm right now.\n", h, m)
}

func getMonday(t time.Time) string {
	if t.Weekday() == time.Monday {
		return SOB(t).Format("2006-01-02T15:04:05Z")
	}
	return getMonday(t.AddDate(0, 0, -1))
}

func getSunday(t time.Time) string {
	if t.Weekday() == time.Sunday {
		return EOB(t).Format("2006-01-02T15:04:05Z")
	}
	return getSunday(t.AddDate(0, 0, 1))
}
