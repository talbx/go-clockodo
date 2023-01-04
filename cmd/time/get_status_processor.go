package time

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/talbx/go-clockodo/cmd/util"
)

type GetStatusProcessor interface{
	CalculateCurrentTime() string
}

type GetStatusProcessorImpl struct {}


func (processor GetStatusProcessorImpl) Process() (string, string) {
	pid := uuid.New().String()
	var entriesRoot string = "v2/entries"
	var currentClock string = "v2/clock"
	var timeSince string = Bod(time.Now()).Format("2006-01-02T15:04:05Z")
	var timeUntil string = time.Now().Format("2006-01-02T15:04:05Z")

	query := fmt.Sprintf("%s?time_since=%s&time_until=%s", entriesRoot, timeSince, timeUntil)
	var repo util.TimeEntriesResponse
	var clo util.ClockResponse
	util.CallApi(query, &repo)
	util.CallApi(currentClock, &clo)

	//fmt.Printf("there is a clock running since %s", clo.Running.RunningSince)
	fulltime := processEntries(&repo, &clo)

	return fmt.Sprintf("Of all completed time entries, you worked a total of %s today!", fulltime), pid
}

func processEntries(resp *util.TimeEntriesResponse, clock *util.ClockResponse) string {
	var totalDuration = 0
	for _, te := range resp.Entries {
		totalDuration+=te.Duration
	}
	h,_ := time.ParseDuration(fmt.Sprintf("%ds",totalDuration))

	hours := int(math.Trunc(float64(h.Minutes()) / 60))
	minutes := int(h.Minutes()) % 60

	// running clock

	t, _ := time.Parse("2006-01-02T15:04:05Z", clock.Running.RunningSince)
	diff := time.Now().Sub(t)

	// sum up
	hours+=int(diff.Hours())
	minutes+=int(diff.Minutes()) % 60
	if minutes % 60 != minutes {
		hours++
		minutes = minutes % 60
	}

	return fmt.Sprintf("%s:%s", AddLeadingZero(hours), AddLeadingZero(minutes))
}

func Bod(t time.Time) time.Time {
    year, month, day := t.Date()
    return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func AddLeadingZero(num int) string {
	if num < 10 {
		return "0" + strconv.Itoa(num)
	}
	return strconv.Itoa(num)
}