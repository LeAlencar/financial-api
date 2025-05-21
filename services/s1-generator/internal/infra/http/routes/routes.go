package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/domain/services"
	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/infra/handlers"
	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/infra/http/middleware"
	"github.com/leandroalencar/banco-dados/shared/models"
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
func SetupRouter(rabbitmq *utils.RabbitMQ, userService *services.UserService) *gin.Engine {
	router := gin.Default()

	// Create handlers
	userHandler := handlers.NewUserHandler(userService, rabbitmq)
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
	}

	return router
}

func handleUserRegistration(rabbitmq *utils.RabbitMQ) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create user model
		user := &models.User{
			Name:      input.Name,
			Email:     input.Email,
			Password:  input.Password, // Note: Password will be hashed by the consumer
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create message with action type
		message := &models.UserMessage{
			Action: models.UserActionCreate,
			User:   user,
		}

		// Send to RabbitMQ
		err := rabbitmq.PublishMessage("users", message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user registration"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message": "User registration request accepted",
			"user": gin.H{
				"name":  user.Name,
				"email": user.Email,
			},
		})
	}
}

func handleUserUpdate(rabbitmq *utils.RabbitMQ) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get ID from URL parameter
		id := c.Param("id")
		var userID int32
		_, err := fmt.Sscanf(id, "%d", &userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}

		// Parse update input
		var input UpdateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create user model with updates
		user := &models.User{
			ID:        userID,
			Name:      input.Name,
			Email:     input.Email,
			Password:  input.Password,
			UpdatedAt: time.Now(),
		}

		// Create message with action type
		message := &models.UserMessage{
			Action: models.UserActionUpdate,
			User:   user,
		}

		// Send to RabbitMQ
		err = rabbitmq.PublishMessage("users", message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user update"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message": "User update request accepted",
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
		})
	}
}

func handleUserDelete(rabbitmq *utils.RabbitMQ) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var userID int32
		_, err := fmt.Sscanf(id, "%d", &userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}

		message := &models.UserMessage{
			Action: models.UserActionDelete,
			User:   &models.User{ID: userID},
		}

		err = rabbitmq.PublishMessage("users", message)

	}
}
