// handlers/tasks.go

package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Task represents the task model
type Task struct {
	Id        int       `json:"Id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	EventID   int       `json:"event_id"`
}

func CreateTask(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var task Task
		var err error

		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// No need to call time.Parse again, already have parsed times

		// Check if event_id exists (no change needed)
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

		// Insert task into database (use time.Time arguments)
		_, err = db.Exec(context.Background(), "INSERT INTO tasks (start_time, end_time, event_id) VALUES ($1, $2, $3)", task.StartTime, task.EndTime, task.EventID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusCreated)
	}
}

func GetTasks(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Define an empty slice to store retrieved tasks
		var tasks []Task

		// Execute SQL query to fetch all tasks
		rows, err := db.Query(context.Background(), "SELECT * FROM tasks")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close() // Close the rows after the function execution

		// Scan each row and append the Task object to the tasks slice
		for rows.Next() {
			var task Task
			err := rows.Scan(&task.Id, &task.StartTime, &task.EndTime, &task.EventID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			tasks = append(tasks, task)
		}

		// Handle any potential errors after rows.Next() loop
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating through rows: " + err.Error()})
			return
		}

		// Respond with the list of tasks in JSON format
		c.JSON(http.StatusOK, tasks)
	}
}

func GetTask(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id") // Get id from URL parameter

		// Convert id to integer and handle potential error
		var id int
		var err error
		id, err = strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		// Define an empty Task struct
		var task Task

		// Query for the task with the given ID
		row := db.QueryRow(context.Background(), "SELECT * FROM tasks WHERE id = $1", id)
		err = row.Scan(&task.Id, &task.StartTime, &task.EndTime, &task.EventID)

		// Check for errors during scan
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		// Respond with the task in JSON format (assuming Task is JSON serializable)
		c.JSON(http.StatusOK, task)
	}
}

func UpdateTask(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id") // Get id from URL parameter

		// Convert id to integer and handle potential error
		var id int
		var err error
		id, err = strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		// Bind the updated task data from JSON request
		var updatedTask Task
		if err := c.BindJSON(&updatedTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if event_id exists (assuming it can't be changed)
		var exists bool
		err = db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)", updatedTask.EventID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
			return
		}

		// Update the task in the database
		result, err := db.Exec(context.Background(), `
		 UPDATE tasks
		 SET start_time = $2, end_time = $3
		 WHERE id = $1
	  `, id, updatedTask.StartTime, updatedTask.EndTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rowsAffected := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		// Respond with No Content (204) on successful update
		c.Status(http.StatusNoContent)
	}
}

func DeleteTask(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id") // Get id from URL parameter

		// Convert id to integer and handle potential error
		var id int
		var err error
		id, err = strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		// Delete task by id
		result, err := db.Exec(context.Background(), "DELETE FROM tasks WHERE id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get the number of rows affected
		rowsAffected := result.RowsAffected() // Only assign the single returned value

		// Respond with No Content (204) if a row was deleted, Not Found (404) otherwise
		if rowsAffected > 0 {
			c.Status(http.StatusNoContent)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		}
	}
}
