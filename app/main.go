package main

import (
	"context"
	"log"

	"sched/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Database connection
	const DATABASE_URL string = "postgres://goatcheese:nHY67ujm@10.0.2.8:5432/schedule"

	db, err := pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	// Initialize Gin
	r := gin.Default()

	// CRUD endpoints for User
	r.POST("/users", handlers.CreateUser(db))
	r.GET("/users", handlers.GetUsers(db))
	r.GET("/users/:id", handlers.GetUser(db))
	r.PUT("/users/:id", handlers.UpdateUser(db))
	r.DELETE("/users/:id", handlers.DeleteUser(db))

	// CRUD endpoints for Event
	r.POST("/events", handlers.CreateEvent(db))
	r.GET("/events", handlers.GetEvents(db))
	r.GET("/events/:id", handlers.GetEvent(db))
	r.PUT("/events/:id", handlers.UpdateEvent(db))
	r.DELETE("/events/:id", handlers.DeleteEvent(db))

	//CRUD endpoints for Task
	r.POST("/tasks", handlers.CreateTask(db))
	r.GET("/tasks", handlers.GetTasks(db))
	r.GET("/tasks/:id", handlers.GetTask(db))
	r.PUT("/tasks/:id", handlers.UpdateTask(db))
	r.DELETE("/tasks/:id", handlers.DeleteTask(db))

	//CRUD endpoints for Routine
	r.POST("/routines", handlers.CreateRoutine(db))
	r.GET("/routines", handlers.GetRoutines(db))
	r.GET("/routines/:id", handlers.GetRoutine(db))
	r.PUT("/routines/:id", handlers.UpdateRoutine(db))
	r.DELETE("/routines/:id", handlers.DeleteRoutine(db))

	// Run the server
	r.Run(":9080")
}
