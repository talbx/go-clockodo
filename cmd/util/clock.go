package util

type ClockResponse struct {
	Running Running `json:"running"`
}
type Running struct {
	RunningSince string `json:"time_since"`
	Id           int64  `json:"id"`
}
