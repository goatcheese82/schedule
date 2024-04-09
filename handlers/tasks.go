// handlers/tasks.go

package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Task represents the task model
type Task struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// CreateTask creates a new task
func CreateTask(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task Task
		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(context.Background(), "INSERT INTO tasks (start_time, end_time) VALUES ($1, $2)", task.StartTime, task.EndTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusCreated)
	}
}

// GetTasks retrieves all tasks
func GetTasks(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(context.Background(), "SELECT start_time, end_time FROM tasks")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var task Task
			if err := rows.Scan(&task.StartTime, &task.EndTime); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			tasks = append(tasks, task)
		}

		c.JSON(http.StatusOK, tasks)
	}
}

// GetTask retrieves a single task by start_time
func GetTask(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := c.Param("start_time")

		var task Task
		err := db.QueryRow(context.Background(), "SELECT start_time, end_time FROM tasks WHERE start_time = $1", startTime).Scan(&task.StartTime, &task.EndTime)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, task)
	}
}

// UpdateTask updates an existing task
func UpdateTask(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := c.Param("start_time")

		var task Task
		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(context.Background(), "UPDATE tasks SET end_time = $1 WHERE start_time = $2", task.EndTime, startTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}

// DeleteTask deletes a task by start_time
func DeleteTask(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := c.Param("start_time")

		_, err := db.Exec(context.Background(), "DELETE FROM tasks WHERE start_time = $1", startTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
