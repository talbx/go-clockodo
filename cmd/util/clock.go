package util
type ClockResponse struct {
	Running Running `json:"running"` 
	/**
{
	"running": {
		"id": 85914058,
		"customers_id": 2350111,
		"projects_id": 2042836,
		"users_id": 234579,
		"billable": 0,
		"texts_id": 30510608,
		"text": "go-clockodo",
		"time_since": "2023-01-03T19:55:11Z",
		"time_until": null,
		"time_insert": "2023-01-03T19:55:11Z",
		"time_last_change": "2023-01-03T19:55:21Z",
		"type": 1,
		"services_id": 860608,
		"duration": null,
		"time_last_change_work_time": "2023-01-03T19:55:11Z",
		"time_clocked_since": "2023-01-03T19:55:11Z",
		"clocked": true,
		"clocked_offline": false,
		"hourly_rate": 84.03
	},
	"current_time": "2023-01-03T19:55:31.742Z"
}
	**/
}

type Running struct {
	RunningSince string `json:"time_since"`
}