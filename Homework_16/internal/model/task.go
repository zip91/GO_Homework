package model

// Task — запись о задаче
type Task struct {
	ID     int64  `json:"id"`
	UID    string `json:"uid"`
	Title  string `json:"title"`
	IsDone bool   `json:"is_done"`
}
