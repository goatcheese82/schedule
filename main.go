package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"sched/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User represents the user model
type User struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Birthday time.Time `json:"birthday"`
}

func main() {
	// Database connection
	db, err := pgxpool.New(context.Background(), "postgres://goatcheese:nHY67ujm@10.0.2.8:5432/schedule")
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	// Initialize Gin
	r := gin.Default()

	// CRUD endpoints for User
	r.POST("/users", handlers.CreateUser(db))
	r.GET("/users", handlers.GetUsers(db))
	r.GET("/users/:username", handlers.GetUser(db))
	r.PUT("/users/:username", handlers.UpdateUser(db))
	r.DELETE("/users/:username", handlers.DeleteUser(db))

	// CRUD endpoints for Event
	r.POST("/events", handlers.CreateEvent(db))
	r.GET("/events", handlers.GetEvents(db))
	r.GET("/events/:title", handlers.GetUser(db))
	r.PUT("/events/:title", handlers.UpdateUser(db))
	r.DELETE("/events/:title", handlers.DeleteUser(db))

	//CRUD endpoints for Task
	r.POST("/tasks", handlers.CreateTask(db))
	r.GET("/tasks", handlers.GetTasks(db))
	r.GET("/tasks/:id", handlers.GetTask(db))
	r.PUT("/tasks/:id", handlers.UpdateTask(db))
	r.DELETE("/tasks/:id", handlers.DeleteTask(db))

	// Run the server
	r.Run(":9080")
}

// User Endpoints

// createUser creates a new user
func createUser(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(context.Background(), "INSERT INTO users (username, email, birthday) VALUES ($1, $2, $3)", user.Username, user.Email, user.Birthday)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusCreated)
	}
}

// getUsers retrieves all users
func getUsers(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(context.Background(), "SELECT username, email, birthday FROM users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var user User
			if err := rows.Scan(&user.Username, &user.Email, &user.Birthday); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, users)
	}
}

// getUser retrieves a single user by username
func getUser(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		var user User
		err := db.QueryRow(context.Background(), "SELECT username, email, birthday FROM users WHERE username = $1", username).Scan(&user.Username, &user.Email, &user.Birthday)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// updateUser updates an existing user
func updateUser(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(context.Background(), "UPDATE users SET email = $1, birthday = $2 WHERE username = $3", user.Email, user.Birthday, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}

// deleteUser deletes a user by username
func deleteUser(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		_, err := db.Exec(context.Background(), "DELETE FROM users WHERE username = $1", username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// Event
