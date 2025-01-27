package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/givko/hoodie/internal/api/contracts"
)

type UserService interface {
	RegisterUser(createUser *contracts.CreateUser) error
	LoginUser(loginUser *contracts.LoginUser) (string, error)
	ListUsers() ([]contracts.ListUser, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterUserHandler is the handler for the POST /api/users/register route.
func (handler *UserHandler) RegisterUserHandler(c *gin.Context) {
	var user contracts.CreateUser

	// Bind JSON body to createDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.service.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (handler *UserHandler) LoginUserHandler(c *gin.Context) {
	var user contracts.LoginUser

	// Bind JSON body to createDto
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := handler.service.LoginUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"jwt": jwt})
}

func (handler *UserHandler) ListUsersHandler(c *gin.Context) {
	users, err := handler.service.ListUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
