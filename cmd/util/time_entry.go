package util

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
	StartTime 	string `json:"time_since"`
}
