package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Routine struct {
	Id          int    `json: "Id"`
	Name        string `json: "name"`
	Description string `json: "description`
}

// CreateRoutine creates a new routine
func CreateRoutine(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var routine Routine
		if err := c.BindJSON(&routine); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(context.Background(), "INSERT INTO routines (name, description) VALUES ($1, $2)", routine.Name, routine.Description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusCreated)
	}
}

// GetRoutines retrieves all routines
func GetRoutines(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(context.Background(), "SELECT id, name, description FROM routines")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var routines []Routine
		for rows.Next() {
			var routine Routine
			if err := rows.Scan(&routine.Id, &routine.Name, &routine.Description); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			routines = append(routines, routine)
		}

		c.JSON(http.StatusOK, routines)
	}
}

// GetRoutine retrieves a single routine by Id
func GetRoutine(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var routine Routine
		err := db.QueryRow(context.Background(), "SELECT id, name, description FROM routines WHERE id = $1", id).Scan(&routine.Id, &routine.Name, &routine.Description)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Routine not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, routine)
	}
}

// UpdateRoutine updates an existing routine
func UpdateRoutine(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var routine Routine
		if err := c.BindJSON(&routine); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec(context.Background(), "UPDATE routines SET name = $1, description = $2 WHERE id = $3", routine.Name, routine.Description, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Routine updated successfully"})
	}
}

// DeleteRoutine deletes a routine by Id
func DeleteRoutine(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec(context.Background(), "DELETE FROM routines WHERE id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
