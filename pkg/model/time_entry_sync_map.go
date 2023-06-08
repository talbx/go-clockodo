package model

import (
	"sync"
)

type TimeEntrySyncMap[T SomeType] struct {
	m     map[string][]T
	mutex sync.Mutex
}

type SomeType interface {
	DayByCustomer | TimeEntry
}

func CreateTE() *TimeEntrySyncMap[TimeEntry] {
	var g TimeEntrySyncMap[TimeEntry]
	g.m = make(map[string][]TimeEntry, 0)
	return &g
}

func CreateDBC() *TimeEntrySyncMap[DayByCustomer] {
	var g TimeEntrySyncMap[DayByCustomer]
	g.m = make(map[string][]DayByCustomer, 0)
	return &g
}

func (m *TimeEntrySyncMap[T]) AppendNonExistent(wd string, entry T) bool {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	m.m[wd] = append(m.m[wd], entry)
	return true
}

func (m *TimeEntrySyncMap[T]) Get() *map[string][]T {
	return &m.m
}
