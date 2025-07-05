package database

import (
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/amblessedezejim/task_manager_api/models"
)

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (repo *TaskRepository) GetTasks(filters models.TaskFilters) ([]models.Task, *models.Pagination, error) {
	var tasks []models.Task
	var totalCount int

	// Create WHERE filter expression
	whereStr := []string{}
	args := []any{}

	if filters.Completed != nil {
		whereStr = append(whereStr, "completed = ?")
		args = append(args, *filters.Completed)
	}

	if filters.Search != "" {
		whereStr = append(whereStr, "(title LIKE ? OR description LIKE ?)")
		searchTerm := "%" + filters.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	var whereSql string
	if len(whereStr) > 0 {
		whereSql = "WHERE " + strings.Join(whereStr, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tasks %s", whereSql)
	err := repo.DB.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(filters.Limit)))
	offset := (filters.Page - 1) * filters.Limit
	args = append(args, filters.Limit, offset)

	pagination := &models.Pagination{
		CurrentPage: filters.Page,
		PerPage:     filters.Limit,
		TotalPages:  totalPages,
		TotalItems:  totalCount,
	}

	query := fmt.Sprintf(`
		SELECT id, title, description, completed, created_at, updated_at
		FROM tasks %s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
		`, whereSql)

	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)

		if err != nil {
			return nil, nil, err
		}

		tasks = append(tasks, task)
	}
	return tasks, pagination, nil
}

func (repo *TaskRepository) GetTaskById(id int) (*models.Task, error) {
	var task models.Task
	query := "SELECT id, title, description, completed, created_at, updated_at FROM tasks WHERE id = ?"
	err := repo.DB.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}

func (repo *TaskRepository) CreateTask(request models.CreateTaskRequest) (*models.Task, error) {
	query := "INSERT INTO tasks (title, description, completed, created_at, updated_at) VALUES (?, ?, false, NOW(), NOW())"
	result, err := repo.DB.Exec(query, request.Title, request.Description)
	if err != nil {
		return nil, err
	}

	taskId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return repo.GetTaskById(int(taskId))
}

func (repo *TaskRepository) UpdateTask(id int, request models.UpdateTaskRequest) (*models.Task, error) {
	existingTask, err := repo.GetTaskById(id)
	if err != nil {
		return nil, err
	}

	if existingTask == nil {
		return nil, nil
	}

	updates := []string{}
	args := []any{}

	if request.Title != nil {
		updates = append(updates, "title = ?")
		args = append(args, *request.Title)
	}

	if request.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *request.Description)
	}

	if request.Completed != nil {
		updates = append(updates, "completed = ?")
		args = append(args, *request.Completed)
	}

	if len(updates) == 0 {
		return existingTask, nil
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = ?", strings.Join(updates, ", "))
	_, err = repo.DB.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return repo.GetTaskById(id)
}

func (repo *TaskRepository) DeleteTask(id int) error {
	query := "DELETE FROM tasks WHERE id = ?"
	result, err := repo.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
