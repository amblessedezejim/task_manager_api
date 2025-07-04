package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/amblessedezejim/task_manager_api/config"
	"github.com/amblessedezejim/task_manager_api/models"
	"github.com/gin-gonic/gin"
)

// Get all tasks from database
func GetTasks(ctx *gin.Context) {
	rows, err := config.DB.Query("SELECT id, title, description, completed, created_at, updated_at FROM tasks ORDER BY created_at DESC")
	if err != nil {
		log.Println("Failed to get rows: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.TaskResponse{
			Status:  "error",
			Message: "Failed to fetch tasks from database",
		})
		return
	}
	defer rows.Close()
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, models.TaskResponse{
				Status:  "error",
				Message: "Failed to scan task",
			})
			return
		}
		tasks = append(tasks, task)
	}
	ctx.JSON(http.StatusOK, models.TaskResponse{
		Status:  "success",
		Message: "Fetched all tasks",
		Data:    tasks,
	})
}

// Fetch task by ID
func GetTaskById(ctx *gin.Context) {
	// Grab task id from parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.TaskResponse{
			Status:  "error",
			Message: "Invalid request ID",
		})
		return
	}

	// Try fetching task from database
	var task models.Task
	query := "SELECT id, title, description, completed, created_at, updated_at FROM tasks WHERE id = ?"
	err = config.DB.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		var msg string
		var code int
		if err == sql.ErrNoRows {
			msg = "Task not found"
			code = http.StatusNotFound
		} else {
			msg = "Failed to fetch task" + err.Error()
			code = http.StatusInternalServerError
		}
		serverError(ctx, msg, code)
		return
	}
	ctx.JSON(http.StatusOK, models.TaskResponse{
		Status:  "success",
		Message: "Task retrieved successfully",
		Data:    task,
	})
}

// Create a new task
func CreateTask(ctx *gin.Context) {
	var request models.CreateTaskRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.TaskResponse{
			Status:  "error",
			Message: "Invalid request",
		})
		return
	}

	query := "INSERT INTO tasks (title description) VALUES (? ?)"
	result, err := config.DB.Exec(query, request.Title, request.Description)
	if err != nil {
		log.Println(err.Error())
		serverError(ctx, fmt.Sprintf("Failed to create task: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	taskId, err := result.LastInsertId()
	if err != nil {
		serverError(ctx, fmt.Sprintf("Failed to get taskID %s", err.Error()), http.StatusInternalServerError)
		return
	}
	var task models.Task
	query = "SELECT id, title, description, completed, created_at, updated_at FROM tasks WHERE id = ?"
	err = config.DB.QueryRow(query, taskId).Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		serverError(ctx, fmt.Sprintf("Fetch to fetch created task: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusCreated, models.TaskResponse{
		Status:  "success",
		Message: "Task created",
		Data:    task,
	})
}

func UpdateTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.TaskResponse{
			Status:  "error",
			Message: "Invalid task ID",
		})
		return
	}

	var request models.UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		serverError(ctx, "Invalid request data "+err.Error(), http.StatusBadRequest)
		return
	}

	// Update task in database
	query := "UPDATE tasks SET title = ?, description = ? , completed = ? WHERE id = ?"
	result, err := config.DB.Exec(query, request.Title, request.Description, request.Completed, id)
	if err != nil {
		serverError(ctx, "Failed to update task "+err.Error(), http.StatusInternalServerError)
		return
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		serverError(ctx, "Failed to get updated results "+err.Error(), http.StatusInternalServerError)
		return
	}

	if affectedRows == 0 {
		serverError(ctx, "Task not found", http.StatusNotFound)
		return
	}

	var task models.Task
	query = "SELECT id, title, description, completed, created_at, updated_at FROM tasks WHERE id = ?"
	err = config.DB.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		serverError(ctx, "Failed to fetch updated task", http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, models.TaskResponse{
		Status:  "success",
		Message: "Task updated successfully",
		Data:    task,
	})
}

// Delete task from database
func DeleteTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.TaskResponse{
			Status:  "error",
			Message: "Invalid task ID",
		})
		return
	}

	var request models.UpdateTaskRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		serverError(ctx, "Invalid request data "+err.Error(), http.StatusBadRequest)
		return
	}

	// Delete task in database
	query := "DELETE FROM tasks WHERE id = ?"
	result, err := config.DB.Exec(query, id)
	if err != nil {
		serverError(ctx, "Failed to delete task "+err.Error(), http.StatusInternalServerError)
		return
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		serverError(ctx, "Failed to get delete results "+err.Error(), http.StatusInternalServerError)
		return
	}

	if affectedRows == 0 {
		serverError(ctx, "Task not found", http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, models.TaskResponse{
		Status:  "success",
		Message: "Task deleted successfully",
	})
}

func serverError(ctx *gin.Context, msg string, statusCode int) {
	ctx.JSON(statusCode, models.TaskResponse{
		Status:  "error",
		Message: msg,
	})

}
