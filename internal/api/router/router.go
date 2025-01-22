package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/plamendelchev/hoodie/internal/api/handlers"
	"github.com/plamendelchev/hoodie/internal/api/ws"
	"github.com/plamendelchev/hoodie/internal/data/in_memory"
	"github.com/plamendelchev/hoodie/internal/security"
	"github.com/plamendelchev/hoodie/internal/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//TODO: Implement a proper check for the origin
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

	router.GET("/ws", wsHandler)
	go Hub.Run()

	return router
}

func setupUsersApiRoutes(router *gin.Engine) {
	user_in_memory_repository := in_memory.UserInMemoryRepository{}
	user_service := service.NewUserService(user_in_memory_repository)
	user_handler := handlers.NewUserHandler(user_service)

	// Example of setting up a group of routes
	api_users := router.Group("/api/users")
	api_users.POST("/register", user_handler.RegisterUserHandler)
	api_users.POST("/login", user_handler.LoginUserHandler)
}

func setupAdminApiRoutes(router *gin.Engine) {
	api_admin := router.Group("/api/admin")
	api_admin.Use(adminOnlyMiddleware())
	api_admin.GET("/users", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "admin route"})
	})

	// You can also set up other route groups or standalone routes
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
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	/*	for {
		message := ws.Message{}
		err := conn.ReadJSON(&message)
		if err != nil {
			break
		}

		conn.WriteJSON(message)
	}*/

	username := c.Query("username")
	wsConn := ws.NewWsConnection(conn, Hub, username)
	Hub.Register <- wsConn
}
