package model

import (
	"sync"
	"time"
)

type TimeEntrySyncMap[T SomeType] struct {
	m     map[time.Weekday][]T
	mutex sync.Mutex
}

type SomeType interface {
	DayByCustomer | TimeEntry
}

func CreateTE() *TimeEntrySyncMap[TimeEntry] {
	var g TimeEntrySyncMap[TimeEntry]
	g.m = make(map[time.Weekday][]TimeEntry, 0)
	return &g
}

func CreateDBC() *TimeEntrySyncMap[DayByCustomer] {
	var g TimeEntrySyncMap[DayByCustomer]
	g.m = make(map[time.Weekday][]DayByCustomer, 0)
	return &g
}

func (m *TimeEntrySyncMap[T]) AppendNonExistent(wd time.Weekday, entry T) bool {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	m.m[wd] = append(m.m[wd], entry)
	return true
}

func (m *TimeEntrySyncMap[T]) Get() *map[time.Weekday][]T {
	return &m.m
}
