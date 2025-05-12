package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/domain/services"
	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/infra/http/handlers"
	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/infra/http/middleware"
)

// SetupRouter configures all the routes for the application
func SetupRouter(userService *services.UserService) *gin.Engine {
	router := gin.Default()

	// Create handlers
	authHandler := handlers.NewAuthHandler(userService)
	userHandler := handlers.NewUserHandler(userService)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		users := protected.Group("/users")
		{
			users.GET("/:id", userHandler.GetUser)
		}

		// TODO: Add other protected routes (transactions, quotations, etc.)
	}

	return router
}
