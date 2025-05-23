package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/domain/services"
	"github.com/leandroalencar/banco-dados/shared/messaging/events"
	"github.com/leandroalencar/banco-dados/shared/utils"
)

type UserHandler struct {
	userService        *services.UserService
	transactionService *services.TransactionService
	rabbitmq           *utils.RabbitMQ
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type TransactionInput struct {
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	CurrencyPair string  `json:"currency_pair" binding:"required"`
	QuotationID  string  `json:"quotation_id,omitempty"`
}

func NewUserHandler(userService *services.UserService, transactionService *services.TransactionService, rabbitmq *utils.RabbitMQ) *UserHandler {
	return &UserHandler{
		userService:        userService,
		transactionService: transactionService,
		rabbitmq:           rabbitmq,
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

func (h *UserHandler) Buy(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create transaction event data
	eventData := events.TransactionEventData{
		UserID:       fmt.Sprintf("%v", userID),
		CurrencyPair: input.CurrencyPair,
		Amount:       input.Amount,
		QuotationID:  input.QuotationID,
		Timestamp:    time.Now(),
	}

	// Create message with action type
	message := &events.TransactionEvent{
		Action: events.TransactionActionBuy,
		Data:   eventData,
	}

	// Send to RabbitMQ
	err := h.rabbitmq.PublishMessage("transactions", message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process buy transaction"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Buy transaction request accepted",
		"transaction": gin.H{
			"user_id":       eventData.UserID,
			"currency_pair": eventData.CurrencyPair,
			"amount":        eventData.Amount,
			"type":          "BUY",
			"timestamp":     eventData.Timestamp,
		},
	})
}

func (h *UserHandler) Sell(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create transaction event data
	eventData := events.TransactionEventData{
		UserID:       fmt.Sprintf("%v", userID),
		CurrencyPair: input.CurrencyPair,
		Amount:       input.Amount,
		QuotationID:  input.QuotationID,
		Timestamp:    time.Now(),
	}

	// Create message with action type
	message := &events.TransactionEvent{
		Action: events.TransactionActionSell,
		Data:   eventData,
	}

	// Send to RabbitMQ
	err := h.rabbitmq.PublishMessage("transactions", message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process sell transaction"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Sell transaction request accepted",
		"transaction": gin.H{
			"user_id":       eventData.UserID,
			"currency_pair": eventData.CurrencyPair,
			"amount":        eventData.Amount,
			"type":          "SELL",
			"timestamp":     eventData.Timestamp,
		},
	})
}

func (h *UserHandler) GetTransactions(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get limit from query parameter (default to 50)
	limitStr := c.DefaultQuery("limit", "50")
	limit := 50
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	// Get transactions from MongoDB
	transactions, err := h.transactionService.GetUserTransactions(c.Request.Context(), fmt.Sprintf("%v", userID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve transactions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"count":        len(transactions),
		"user_id":      fmt.Sprintf("%v", userID),
		"limit":        limit,
	})
}
