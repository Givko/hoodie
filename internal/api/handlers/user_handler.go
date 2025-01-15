package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plamendelchev/hoodie/internal/api/dto"
)

type UserService interface {
	RegisterUser(createUser dto.CreateUser) error
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterUserHandler is the handler for the POST /api/users/register route.
func (handler *UserHandler) RegisterUserHandler(c *gin.Context) {
	var user dto.CreateUser

	// Bind JSON body to createDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.service.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
