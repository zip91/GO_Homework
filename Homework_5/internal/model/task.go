package model

type Task struct {
	ID     int    `json:"id"`
	UID    string `json:"uid"`
	Title  string `json:"title"`
	IsDone bool   `json:"is_done"`
}
