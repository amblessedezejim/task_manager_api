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
	ID          int       `json:"id" db:"id"`
	Description string    `json:"description" db:"description"`
	Title       string    `json:"title" db:"title"`
	Completed   bool      `json:"completed" db:"completed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type TaskResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type CreateTaskRequest struct {
	Title       string `json:"status"`
	Description string `json:"message"`
}

type UpdateTaskRequest struct {
	Title       string `json:"status"`
	Description string `json:"message"`
	Completed   bool   `json:"completed"`
}
