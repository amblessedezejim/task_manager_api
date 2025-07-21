package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/amblessedezejim/task_manager_api/database"
	"github.com/amblessedezejim/task_manager_api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TaskHandler struct {
	repo     *database.TaskRepository
	validate *validator.Validate
}

func NewTaskHandler(db *sql.DB) *TaskHandler {
	return &TaskHandler{
		repo:     database.NewTaskRepository(db),
		validate: validator.New(),
	}
}

// Get all tasks from database
func (handler *TaskHandler) GetTasks(ctx *gin.Context) {
	var filters models.TaskFilters

	if err := ctx.ShouldBindQuery(&filters); err != nil {
		handler.errorResponse(ctx, "Invalid query parameters", http.StatusBadRequest, err)
		return
	}

	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 || filters.Limit > 100 {
		filters.Limit = 10
	}

	tasks, pagination, err := handler.repo.GetTasks(filters)
	if err != nil {
		handler.errorResponse(ctx, "Failed to retreive tasks", http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, models.TaskListResponse{
		Status:     "success",
		Message:    "Retrieved all tasks",
		Data:       tasks,
		Pagination: pagination,
	})
}

// Fetch task by ID
func (handler *TaskHandler) GetTaskById(ctx *gin.Context) {
	// Grab task id from parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		handler.errorResponse(ctx, "Invalid task ID", http.StatusBadRequest, err)
		return
	}

	// Try fetching task from database
	task, err := handler.repo.GetTaskById(id)
	if err != nil {
		handler.errorResponse(ctx, "Failed to retreive task", http.StatusInternalServerError, err)
	}

	if task == nil {
		handler.errorResponse(ctx, "Task not found", http.StatusNotFound, nil)
	}
	ctx.JSON(http.StatusOK, models.TaskResponse{
		Status:  "success",
		Message: "Task retrieved successfully",
		Data:    task,
	})
}

// Create a new task
func (handler *TaskHandler) CreateTask(ctx *gin.Context) {
	var request models.CreateTaskRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		handler.validationErrorResponse(ctx, err)
		return
	}

	task, err := handler.repo.CreateTask(request)
	if err != nil {
		handler.errorResponse(ctx, "Failed to create task", http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, models.TaskResponse{
		Status:  "success",
		Message: "Task created successfully",
		Data:    task,
	})
}

// Updates task in database
func (handler *TaskHandler) UpdateTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		handler.validationErrorResponse(ctx, err)
		return
	}

	var request models.UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		handler.validationErrorResponse(ctx, err)
		return
	}

	// TODO:
	// if err := handler.validate.Struct(request); err != nil {
	// 	handler.validationErrorResponse(ctx, err)
	// 	return
	// }

	task, err := handler.repo.UpdateTask(id, request)
	if err != nil {
		handler.errorResponse(ctx, "Failed to update task", http.StatusInternalServerError, err)
		return
	}

	if task == nil {
		handler.errorResponse(ctx, "Task not found", http.StatusNotFound, nil)
	}

	ctx.JSON(http.StatusOK, models.TaskResponse{
		Status:  "success",
		Message: "Task updated successfully",
		Data:    task,
	})
}

// Delete task from database
func (handler *TaskHandler) DeleteTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		handler.errorResponse(ctx, "Invalid task ID", http.StatusBadRequest, err)
		return
	}

	err = handler.repo.DeleteTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			handler.errorResponse(ctx, "Task not found", http.StatusNotFound, nil)
			return
		}
		handler.errorResponse(ctx, "Failed to delete task", http.StatusInternalServerError, err)
		return

	}

	ctx.JSON(http.StatusOK, models.TaskResponse{
		Status:  "success",
		Message: "Task deleted successfully",
	})
}

func (handler *TaskHandler) errorResponse(ctx *gin.Context, msg string, statusCode int, err error) {
	response := models.ErrorResponse{
		Status:  "error",
		Message: msg,
	}

	if err != nil && gin.Mode() == gin.DebugMode {
		response.Errors = map[string]any{
			"detail": err.Error(),
		}
	}
	ctx.JSON(statusCode, response)
}

func (handler *TaskHandler) validationErrorResponse(ctx *gin.Context, err error) {
	errors := make(map[string]any)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors[e.Field()] = handler.getValidationMessage(e)
		}
	} else {
		errors["general"] = err.Error()
	}

	ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
		Status:  "error",
		Message: "Validation failed",
		Errors:  errors,
	})
}

func (handler *TaskHandler) getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "This field must be at least " + e.Param() + " characters long"
	case "max":
		return "This field must be at most" + e.Param() + " characters long"
	default:
		return "This field is invalid"
	}
}
