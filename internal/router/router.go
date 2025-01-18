package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/plamendelchev/hoodie/internal/api/handlers"
	"github.com/plamendelchev/hoodie/internal/data/in_memory"
	"github.com/plamendelchev/hoodie/internal/security"
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
	api_users := router.Group("/api/users")
	api_users.POST("/register", user_handler.RegisterUserHandler)
	api_users.POST("/login", user_handler.LoginUserHandler)

	api_admin := router.Group("/api/admin")
	api_admin.Use(AdminOnly())
	api_admin.GET("/users", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "admin route"})
	})

	// You can also set up other route groups or standalone routes
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {

		auth_header := c.GetHeader("Authorization")
		auth_header = strings.TrimPrefix(auth_header, "Bearer ")
		if auth_header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		token, err := security.VerifyToken(auth_header)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		isAdmin := claims["admin"].(bool)

		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}
