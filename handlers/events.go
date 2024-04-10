// handlers/events.go

package handlers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Event represents the event model
type Event struct {
	Id    int    `json:"Id"`
	Title string `json:"title"`
	Image string `json:"image"`
}

// CreateEvent creates a new event
func CreateEvent(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event Event
		if err := c.BindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(context.Background(), "INSERT INTO events (title, image) VALUES ($1, $2)", event.Title, event.Image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusCreated)
	}
}

// GetEvents retrieves all events
func GetEvents(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(context.Background(), "SELECT title, image FROM events")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var events []Event
		for rows.Next() {
			var event Event
			if err := rows.Scan(&event.Title, &event.Image); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			events = append(events, event)
		}

		c.JSON(http.StatusOK, events)
	}
}

// GetEvent retrieves a single event by title
func GetEvent(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")

		var event Event
		err := db.QueryRow(context.Background(), "SELECT title, image FROM events WHERE title = $1", title).Scan(&event.Title, &event.Image)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, event)
	}
}

// UpdateEvent updates an existing event
func UpdateEvent(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")

		var event Event
		if err := c.BindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(context.Background(), "UPDATE events SET image = $1 WHERE title = $2", event.Image, title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	}
}

// DeleteEvent deletes an event by title
func DeleteEvent(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")

		_, err := db.Exec(context.Background(), "DELETE FROM events WHERE title = $1", title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
