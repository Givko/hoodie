package router

import (
	"github.com/gin-gonic/gin"
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
	// Example of setting up a group of routes
	api := router.Group("/api")
	{
		api.POST("/users/register", func(c *gin.Context) {})
		// Add more routes as needed
	}

	// You can also set up other route groups or standalone routes
}
