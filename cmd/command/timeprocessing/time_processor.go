package timeprocessing

import (
	"fmt"
	"sort"
	"time"

	"github.com/talbx/go-clockodo/pkg/concurrent"
	. "github.com/talbx/go-clockodo/pkg/model"
	"github.com/talbx/go-clockodo/pkg/render"
	"github.com/talbx/go-clockodo/pkg/util"
	. "github.com/talbx/go-clockodo/pkg/util"

	"golang.org/x/exp/slices"
)

type TimeProcessor struct{}

func (p TimeProcessor) Process(mode string, last int) error {

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

	for _, v := range *result {
		sort.Slice(v[:], func(i, j int) bool {
			return v[i].Customer < v[j].Customer
		})
	}

	sortedMappy := sortWeekdayMap(*result)
	render.Render(sortedMappy, clo, util.ProcessClock)
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

	sort.Slice(keys, func(i, j int) bool {

		prev, err := time.Parse("02.01.2006", keys[i])
		if err != nil {
			util.SugaredLogger.Error("error parsing date %v", keys[i], err)
		}
		curr, err := time.Parse("02.01.2006", keys[j])
		if err != nil {
			util.SugaredLogger.Error("error parsing date %v", keys[j], err)
		}

		return prev.Before(curr)
	})

	for _, k := range keys {
		r[k] = mappy[k]
	}

	return r
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
