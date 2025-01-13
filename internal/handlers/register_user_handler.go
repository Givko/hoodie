package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/plamendelchev/hoodie/internal/data/models"
)

var users = sync.Map{}

// RegisterUserHandler is the handler for the POST /api/users/register route.
func RegisterUserHandler(c *gin.Context) {
	var user models.User

	// Bind JSON body to createDto
	if err := c.ShouldBindJSON(&user); err != nil {
		// If binding fails, return a 400 Bad Request with the error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a unique ID for the user
	user.Id = uuid.NewString()
	value, isLoaded := users.LoadOrStore(user.Username, user)
	if isLoaded {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    value,
	})
}
