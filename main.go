package main

import (
	"context"
	"log"
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
