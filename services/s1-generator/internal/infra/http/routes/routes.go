package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/domain/services"
	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/infra/handlers"
	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/infra/http/middleware"
	"github.com/leandroalencar/banco-dados/shared/utils"
)

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserInput struct {
	ID       int32  `json:"id" binding:"required"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// SetupRouter configures all the routes for the application
func SetupRouter(rabbitmq *utils.RabbitMQ, userService *services.UserService, transactionService *services.TransactionService) *gin.Engine {
	router := gin.Default()

	// Create handlers
	userHandler := handlers.NewUserHandler(userService, transactionService, rabbitmq)
	authHandler := handlers.NewAuthHandler(userService)
	quotationHandler := handlers.NewQuotationHandler(rabbitmq)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	router.POST("/quotations/generate", quotationHandler.GenerateQuotations)

	users := router.Group("/users")
	{
		users.POST("/register", userHandler.Register)
		users.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Protected user routes
		protectedUsers := protected.Group("/users")
		{
			protectedUsers.GET("/:id", userHandler.GetUser)
			protectedUsers.PATCH("/:id", userHandler.Update)
			protectedUsers.DELETE("/:id", userHandler.Delete)
		}
		protectedTransactions := protected.Group("/transactions")
		{
			protectedTransactions.GET("/", userHandler.GetTransactions)
			protectedTransactions.POST("/buy", userHandler.Buy)
			protectedTransactions.POST("/sell", userHandler.Sell)

		}
	}

	return router
}
