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
	EventID   int       `json:"event_id"` // Foreign key referencing Event ID
}

// // CreateTask creates a new task
// func CreateTask(db *pgxpool.Pool) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var task Task

// 		if err := c.BindJSON(&task); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// Parse start_time using the custom format
// 		startTime, err := time.Parse("15:04:05", task.StartTime.Format("15:04:05"))
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// Parse end_time using the custom format
// 		endTime, err := time.Parse("15:04:05", task.EndTime.Format("15:04:05"))
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		// Set the parsed times back to the Task struct
// 		task.StartTime = startTime
// 		task.EndTime = endTime

// 		_, err = db.Exec(context.Background(), "INSERT INTO tasks (start_time, end_time, event_id) VALUES ($1, $2, $3)", task.StartTime, task.EndTime, task.EventID)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}

// 		c.Status(http.StatusCreated)
// 	}
// }

func CreateTask(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task Task

		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse location from request (optional, depending on your implementation)
		location, err := time.LoadLocation("America/Los_Angeles") // Replace with appropriate location if needed
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse received times (assuming format remains "15:04:05")
		startTime, err := time.ParseInLocation("15:04:05", task.StartTime.Format("15:04:05"), location)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		endTime, err := time.ParseInLocation("15:04:05", task.EndTime.Format("15:04:05"), location)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert parsed times to UTC
		startTime = startTime.In(time.UTC)
		endTime = endTime.In(time.UTC)

		// Check if event_id exists
		var exists bool
		err = db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)", task.EventID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
			return
		}

		// Insert task into database
		_, err = db.Exec(context.Background(), "INSERT INTO tasks (start_time, end_time, event_id) VALUES ($1, $2, $3)", startTime, endTime, task.EventID)
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
