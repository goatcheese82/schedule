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
	db, err := pgxpool.New(context.Background(), "postgres://goatcheese:B8qAieedivIiOiTdtW0VgWWohShfaMSm@dpg-coasek21hbls73ftkiig-a.oregon-postgres.render.com/routines_7w39")
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
