package timeprocessing

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/talbx/go-clockodo/cmd/util"
)

type TimeProcessor interface {
	Process(round bool)
}

type TimeProcessorFactory struct{}

func (t TimeProcessorFactory) CreateTimeProcessor() TimeProcessor {
	return WeekProcessor{}
}

var instance *TimeProcessorFactory

func CreateCommandFactory() *TimeProcessorFactory {
	if instance == nil {
		instance = &TimeProcessorFactory{}
	}
	return instance
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

func DurationToHM(duration int)(hours int, minutes int){
	h, _ := time.ParseDuration(fmt.Sprintf("%ds", duration))

	if h.Minutes() > 0 {
		hours = int(math.Trunc(float64(h.Minutes()) / 60))
		minutes = int(h.Minutes()) % 60
	}
	return
}

func CreateOutput(hours int, minutes int, tasks *string, round bool, today bool) {
	var roundMsg string = ""
	if round {
		hours, minutes = Round(hours, minutes)
		roundMsg = "The time is rounded on 15m basis."
	}
	fulltime := fmt.Sprintf("%s:%s", AddLeadingZero(hours), AddLeadingZero(minutes))

	if today {
		timeMsg := "You worked %s today! %s"
		log.Printf(timeMsg, fulltime, roundMsg)
	} else {
		timeMsg := "You worked %s yesterday! %s"
		log.Printf(timeMsg, fulltime, roundMsg)
	}

	if tasks != nil {
		log.Printf("Your tasks during that time were %s\n", *tasks)
	}
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
		return logTimes(h, m, h, m)
	}
	if r < 8 {
		return logTimes(h, m-r, h, m)
	}
	x := m + (15 - r)
	if x == 60 {
		return logTimes(h+1, 0, h, m)
	}
	return logTimes(h, x, h, m)
}

func logTimes(h int, m int, ih int, im int) (int, int) {
	//fmt.Printf("Rounded Input %d:%d to %d:%d\n", ih, im, h, m)
	return h, m
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value && item != "" && item != " " {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
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
