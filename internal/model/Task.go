package model

type TaskStatus string

const (
	Pending    TaskStatus = "PENDING"
	InProgress TaskStatus = "IN_PROGRESS"
	Done       TaskStatus = "DONE"
)

type Task struct {
	Id      string
	Status  TaskStatus
	Payload string
}
