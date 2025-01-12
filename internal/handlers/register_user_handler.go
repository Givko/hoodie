package handlers

import (
	"net/http"

	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var users = sync.Map{}

type User struct {
	Id       string
	Username string
	password string
}

// RegisterUserHandler is the handler for the POST /api/users/register route.
func RegisterUserHandler(c *gin.Context) {
	var user User

	// Bind JSON body to createDto
	if err := c.ShouldBindJSON(&user); err != nil {
		// If binding fails, return a 400 Bad Request with the error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a unique ID for the user
	user.Id = uuid.NewString()
	value, loaded := users.LoadOrStore(user.Username, user)
	if loaded {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
		"data":    value,
	})
}
