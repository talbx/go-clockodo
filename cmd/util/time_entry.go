package util

type TimeEntriesResponse struct {
	Entries []TimeEntry `json:"entries""`
	/**
	{
	"paging": {
		"items_per_page": 2500,
		"current_page": 1,
		"count_pages": 1,
		"count_items": 6
	},
	"filter": null,
	"entries": [
		**/
}

type TimeEntry struct {
	Id         int `json:"id"`
	CustomerId int `json:"customers_id"`
	ProjectId  int `json:"projects_id"`
	UserId     int `json:"users_id"`
	Billable   int `json:"billable"`
	Duration   int `json:"duration"`

	/**
	{
			"id": 85722133,
			"customers_id": 2410335,
			"projects_id": 2020072,
			"users_id": 234579,
			"billable": 1,
			"texts_id": null,
			"text": null,
			"time_since": "2022-12-30T07:21:04Z",
			"time_until": "2022-12-30T08:45:34Z",
			"time_insert": "2022-12-30T07:21:04Z",
			"time_last_change": "2022-12-30T08:45:34Z",
			"type": 1,
			"services_id": 860608,
			"duration": 5070,
			"time_last_change_work_time": "2022-12-30T08:45:34Z",
			"time_clocked_since": "2022-12-30T07:21:04Z",
			"clocked": true,
			"clocked_offline": false,
			"hourly_rate": 84.03
		},
	**/
}
