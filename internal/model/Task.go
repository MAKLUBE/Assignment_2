package model

type TaskStatus string

const (
	Pending    TaskStatus = "PENDING"
	InProgress TaskStatus = "IN_PROGRESS"
	Done       TaskStatus = "DONE"
)

type Task struct {
	id      int
	status  TaskStatus
	payload string
}
