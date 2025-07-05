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
	Title       string    `json:"title" db:"title" binding:"required,min=1,max=255"`
	Description string    `json:"description" db:"description" validate:"max=1000"`
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
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"max=1000"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty" binding:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=1000"`
	Completed   *bool   `json:"completed,omitempty"`
}

type TaskListResponse struct {
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       any         `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
	TotalItems  int `json:"total_items"`
}

type TaskFilters struct {
	Completed *bool  `form:"completed"`
	Search    string `form:"search"`
	Page      int    `form:"page,default=1"`
	Limit     int    `form:"limit,default=10"`
}

type ErrorResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Errors  map[string]any `json:"errors,omitempty"`
}
