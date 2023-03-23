package model

type Customer struct {
	CustomerName `json:"customer"`
}

type CustomerName struct {
	Name string `json:"name"`
}

type StartPayload struct {
	CustomerId  int    `json:"customers_id"`
	ServiceId   int    `json:"services_id"`
	ProjectId   int    `json:"projects_id"`
	Description string `json:"text"`
}

type TimeEntriesResponse struct {
	Entries []TimeEntry `json:"entries"`
}

type TimeEntry struct {
	Id          int    `json:"id"`
	CustomerId  int    `json:"customers_id"`
	ProjectId   int    `json:"projects_id"`
	UserId      int    `json:"users_id"`
	Billable    int    `json:"billable"`
	Duration    int    `json:"duration"`
	Description string `json:"text"`
	StartTime   string `json:"time_since"`
}

type ClockResponse struct {
	Running Running `json:"running"`
}
type Running struct {
	RunningSince string `json:"time_since"`
	Id           int64  `json:"id"`
}

type DayByCustomer struct {
	Customer        string
	Tasks           string
	TotalTime       int
	AggregatedTime  string
	CustomerId      int
	AggregatedTasks string
}

type GoClockodoConfig struct {
	ApiKey  string
	ApiUser string
}
