package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/givko/hoodie/internal/api/handlers"
	"github.com/givko/hoodie/internal/api/ws"
	"github.com/givko/hoodie/internal/data/in_memory"
	"github.com/givko/hoodie/internal/security"
	"github.com/givko/hoodie/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var Hub = ws.NewHub()

// Init initializes the Gin router with all routes and middleware.
func Init() *gin.Engine {
	router := gin.Default()

	// Initialize routes
	setupUsersApiRoutes(router)
	setupAdminApiRoutes(router)
	setupWebsocketRoutes(router)

	go Hub.Run()

	return router
}

func setupWebsocketRoutes(router *gin.Engine) {
	ws := router.Group("/ws")
	ws.Use(isLoggedMiddleware())
	ws.GET("/connect", wsHandler)
}

func setupUsersApiRoutes(router *gin.Engine) {
	user_in_memory_repository := in_memory.UserInMemoryRepository{}
	user_service := service.NewUserService(user_in_memory_repository)
	user_handler := handlers.NewUserHandler(user_service)

	// Example of setting up a group of routes
	api_users := router.Group("/api/users")
	api_users.POST("/register", user_handler.RegisterUserHandler)
	api_users.POST("/login", user_handler.LoginUserHandler)

	// User list route subgroup requires the user to be logged in
	api_users_list := api_users.Group("/list")
	api_users_list.Use(isLoggedMiddleware())
	api_users_list.GET("/all", user_handler.ListUsersHandler)
}

func setupAdminApiRoutes(router *gin.Engine) {
	api_admin := router.Group("/api/admin")
	api_admin.Use(adminOnlyMiddleware())
	api_admin.GET("/users", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "admin route"})
	})

	// You can also set up other route groups or standalone routes
}

func isLoggedMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		auth_header := c.GetHeader("Authorization")
		auth_header = strings.TrimPrefix(auth_header, "Bearer ")
		if auth_header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "JWT token is missing for Authorization header"})
			return
		}

		token, err := security.VerifyToken(auth_header)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT Token"})
			return
		}

		c.Set("username", claims["sub"])
		c.Next()
	}
}

func adminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth_header := c.GetHeader("Authorization")
		auth_header = strings.TrimPrefix(auth_header, "Bearer ")
		if auth_header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "JWT token is missing for Authorization header"})
			return
		}

		token, err := security.VerifyToken(auth_header)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT Token"})
			return
		}

		isAdmin := claims["admin"].(bool)
		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User is not an admin"})
			return
		}

		c.Next()
	}
}

func wsHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while parsing username"})
		return
	}

	wsConn := ws.NewWsConnection(conn, Hub, usernameStr)
	Hub.Register <- wsConn
}
