package util

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type DayByCustomer struct {
	Customer        string
	Tasks           string
	TotalTime       int
	AggregatedTime  string
	CustomerId      int
	AggregatedTasks string
}

func GroupDayEntriesByCustomer(mappy map[time.Weekday][]TimeEntry, custIds []int) map[time.Weekday][]DayByCustomer {
	mappyNew := make(map[time.Weekday][]DayByCustomer)
	for _, customerId := range custIds {
		for w, v := range mappy {
			dbc := DayByCustomer{CustomerId: customerId}
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

func GroupEntriesByDay(repo TimeEntriesResponse) map[time.Weekday][]TimeEntry {
	weekdayMap := make(map[time.Weekday][]TimeEntry)
	weekdays := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}
	log.Printf("We have a total of %v entries\n", len(repo.Entries))
	var count int = 0
	for _, te := range repo.Entries {
		for _, weekday := range weekdays {
			startTime, _ := time.Parse("2006-01-02T15:04:05Z", te.StartTime)
			count++
			if weekday == startTime.Weekday() {
				//fmt.Printf("Current weekday: %v, weekday of entry: %v, time of entry %v, entry: %+v\n", weekday.String(), startTime.Weekday().String(), startTime, te)
				weekdayMap[weekday] = append(weekdayMap[weekday], te)
				break;
			}
		}
	}
	log.Printf("we have a total loop count of %v\n", count)
	return weekdayMap
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