package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/amblessedezejim/task_manager_api/models"
	"github.com/gin-gonic/gin"
)

func GetTasks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, models.TaskResponse{
		Status:  "success",
		Message: "Return all tasks",
		Data:    models.Tasks,
	})
}

func GetTaskById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.TaskResponse{
			Status:  "error",
			Message: "Invalid request ID",
		})
		return
	}

	for _, task := range models.Tasks {
		if task.ID == id {
			ctx.JSON(http.StatusOK, models.TaskResponse{
				Status:  "success",
				Message: "Fetch all tasks",
				Data:    models.Tasks,
			})
			return
		}
	}
	ctx.JSON(http.StatusNotFound, models.TaskResponse{
		Status:  "error",
		Message: "Task not found",
	})
}

func CreateTask(ctx *gin.Context) {
	var task models.Task

	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, models.TaskResponse{
			Status:  "error",
			Message: "Invalid request",
		})
		return
	}

	task.ID = models.NewId()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	models.Tasks = append(models.Tasks, task)
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

	var newTask models.Task
	if err := ctx.ShouldBindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, models.TaskResponse{
			Status:  "error",
			Message: "Invalid request data",
		})
		return
	}

	for i, task := range models.Tasks {
		if task.ID == id {
			newTask.UpdatedAt = time.Now()
			models.Tasks[i] = newTask

			ctx.JSON(http.StatusOK, models.TaskResponse{
				Status:  "success",
				Message: "Tasks updated",
				Data:    newTask,
			})
		}
	}
}

func DeleteTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.TaskResponse{
			Status:  "error",
			Message: "Invalid task ID",
		})
		return
	}

	for i, task := range models.Tasks {
		if task.ID == id {
			models.Tasks = append(models.Tasks[:1], models.Tasks[i+1:]...)

			ctx.JSON(http.StatusOK, models.TaskResponse{
				Status:  "success",
				Message: "Task deleted",
			})

			return
		}
	}

	ctx.JSON(http.StatusNotFound, models.TaskResponse{
		Status:  "error",
		Message: "Task not found",
	})
}
