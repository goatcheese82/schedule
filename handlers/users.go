// handlers/users.go

package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User represents the user model
type User struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Birthday time.Time `json:"birthday"`
}

// CreateUser creates a new user
func CreateUser(db *pgxpool.Pool) gin.HandlerFunc {
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

// GetUsers retrieves all users
func GetUsers(db *pgxpool.Pool) gin.HandlerFunc {
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

// GetUser retrieves a single user by username
func GetUser(db *pgxpool.Pool) gin.HandlerFunc {
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

// UpdateUser updates an existing user
func UpdateUser(db *pgxpool.Pool) gin.HandlerFunc {
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

// DeleteUser deletes a user by username
func DeleteUser(db *pgxpool.Pool) gin.HandlerFunc {
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
