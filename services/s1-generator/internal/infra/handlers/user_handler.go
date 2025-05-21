package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/domain/services"
	"github.com/leandroalencar/banco-dados/shared/messaging/events"
	"github.com/leandroalencar/banco-dados/shared/utils"
)

type UserHandler struct {
	userService *services.UserService
	rabbitmq    *utils.RabbitMQ
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func NewUserHandler(userService *services.UserService, rabbitmq *utils.RabbitMQ) *UserHandler {
	return &UserHandler{
		userService: userService,
		rabbitmq:    rabbitmq,
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	// Convert string ID to uint
	var userID uint
	_, err := fmt.Sscanf(id, "%d", &userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Register(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create event data
	eventData := events.UserEventData{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	// Create message with action type
	message := &events.UserEvent{
		Action: events.UserActionCreate,
		Data:   eventData,
	}

	// Send to RabbitMQ
	err := h.rabbitmq.PublishMessage("users", message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user registration"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "User registration request accepted",
		"user": gin.H{
			"name":  input.Name,
			"email": input.Email,
		},
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	// Get ID from URL parameter
	id := c.Param("id")
	var userID int32
	_, err := fmt.Sscanf(id, "%d", &userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse update input
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create event data
	eventData := events.UserEventData{
		ID:       userID,
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	// Create message with action type
	message := &events.UserEvent{
		Action: events.UserActionUpdate,
		Data:   eventData,
	}

	// Send to RabbitMQ
	err = h.rabbitmq.PublishMessage("users", message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user update"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "User update request accepted",
		"user": gin.H{
			"id":    userID,
			"name":  input.Name,
			"email": input.Email,
		},
	})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var userID int32
	_, err := fmt.Sscanf(id, "%d", &userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Create message with action type
	message := &events.UserEvent{
		Action: events.UserActionDelete,
		Data: events.UserEventData{
			ID: userID,
		},
	}

	// Send to RabbitMQ
	err = h.rabbitmq.PublishMessage("users", message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user deletion"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "User deletion request accepted",
		"user_id": userID,
	})
}
