package models

import "time"

var Tasks []Task
var taskId int = 0

func NewId() (newId int) {
	newId = taskId
	taskId++
	return
}

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Title       string    `json:"title"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
