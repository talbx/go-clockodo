package timeprocessing

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/talbx/go-clockodo/cmd/util"
)

type TimeProcessor interface {
	Process(round bool)
}

func SOB(t time.Time) time.Time {
	return d(t, 0, 0)
}

func EOB(t time.Time) time.Time {
	return d(t, 23, 59)
}

func d(t time.Time, h int, m int) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, h, m, 0, 0, t.Location())
}

func DurationToHM(duration int) (hours int, minutes int) {
	h, _ := time.ParseDuration(fmt.Sprintf("%ds", duration))

	if h.Minutes() > 0 {
		hours = int(math.Trunc(float64(h.Minutes()) / 60))
		minutes = int(h.Minutes()) % 60
	}
	return
}

func AddLeadingZero(num int) string {
	if num < 10 {
		return "0" + strconv.Itoa(num)
	}
	return strconv.Itoa(num)
}

func Round(h int, m int) (int, int) {
	r := m % 15
	if r == 0 {
		return h, m
	}
	if r < 8 {
		return h, m - r
	}
	x := m + (15 - r)
	if x == 60 {
		return h + 1, 0
	}
	return h, m
}

func ProcessClock(clock *util.ClockResponse) (hours int, minutes int) {
	// running clock
	if clock.Running.RunningSince != "" && clock.Running.RunningSince > "0" {
		t, _ := time.Parse("2006-01-02T15:04:05Z", clock.Running.RunningSince)
		diff := time.Now().Sub(t)
		hours = int(diff.Hours())
		minutes = int(diff.Minutes()) % 60
	}
	return
}
