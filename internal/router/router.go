package router

import (
	"github.com/gin-gonic/gin"
	"github.com/plamendelchev/hoodie/internal/api/handlers"
	"github.com/plamendelchev/hoodie/internal/data/in_memory"
	"github.com/plamendelchev/hoodie/internal/service"
)

// Init initializes the Gin router with all routes and middleware.
func Init() *gin.Engine {
	router := gin.Default()

	// Apply global middleware
	// router.Use(middleware.Logger())
	// router.Use(middleware.Recovery())

	// Initialize routes
	setupRoutes(router)

	return router
}

func setupRoutes(router *gin.Engine) {
	user_in_memory_repository := in_memory.UserInMemoryRepository{}
	user_service := service.NewUserService(user_in_memory_repository)
	user_handler := handlers.NewUserHandler(user_service)

	// Example of setting up a group of routes
	api := router.Group("/api")
	api.POST("/users/register", user_handler.RegisterUserHandler)
	api.POST("/users/login", user_handler.LoginUserHandler)
	// Add more routes as needed

	// You can also set up other route groups or standalone routes
}
