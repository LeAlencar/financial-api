package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/domain/repositories"
	"github.com/leandroalencar/banco-dados/shared/messaging/events"
	"github.com/leandroalencar/banco-dados/shared/models"
	"github.com/leandroalencar/banco-dados/shared/utils"
)

// TransactionService handles business logic related to transactions
type TransactionService struct {
	userRepo                *repositories.UserRepository
	quotationRepo           *repositories.QuotationRepository
	currencyTransactionRepo *repositories.CurrencyTransactionRepository
	rabbitmq                *utils.RabbitMQ
}

// NewTransactionService creates a new transaction service
func NewTransactionService(
	userRepo *repositories.UserRepository,
	quotationRepo *repositories.QuotationRepository,
	currencyTransactionRepo *repositories.CurrencyTransactionRepository,
	rabbitmq *utils.RabbitMQ,
) *TransactionService {
	return &TransactionService{
		userRepo:                userRepo,
		quotationRepo:           quotationRepo,
		currencyTransactionRepo: currencyTransactionRepo,
		rabbitmq:                rabbitmq,
	}
}

// ProcessTransactionEvent handles different types of transaction events
func (s *TransactionService) ProcessTransactionEvent(ctx context.Context, event *events.TransactionEvent) error {
	switch event.Action {
	case events.TransactionActionBuy:
		return s.handleBuy(ctx, event.Data)
	case events.TransactionActionSell:
		return s.handleSell(ctx, event.Data)
	default:
		return fmt.Errorf("unknown transaction action type: %s", event.Action)
	}
}

func (s *TransactionService) handleBuy(ctx context.Context, data events.TransactionEventData) error {
	// Convert user ID to uint
	userID, err := strconv.ParseUint(data.UserID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	// Get user to check balance and update after transaction
	user, err := s.userRepo.GetByID(ctx, uint(userID))
	if err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}

	// Get latest quotation for the currency pair
	quotation, err := s.quotationRepo.GetLatestByCurrencyPair(ctx, data.CurrencyPair)
	if err != nil {
		// If no quotation found, use fallback rates
		return s.handleBuyWithFallbackRate(ctx, data, user)
	}

	// Calculate total cost using real quotation
	totalCost := data.Amount * quotation.BuyPrice

	// Check if user has sufficient balance
	if user.Balance < totalCost {
		// Create and save failed transaction
		failedTransaction := &models.Transaction{
			ID:           generateTransactionID(),
			UserID:       data.UserID,
			Type:         models.Buy,
			CurrencyPair: data.CurrencyPair,
			Amount:       data.Amount,
			ExchangeRate: quotation.BuyPrice,
			TotalValue:   totalCost,
			Status:       "failed: insufficient_balance",
			Timestamp:    data.Timestamp,
			QuotationID:  quotation.ID,
		}

		// Save failed transaction to MongoDB
		if err := s.currencyTransactionRepo.Create(ctx, failedTransaction); err != nil {
			return fmt.Errorf("failed to save failed transaction: %v", err)
		}

		// Send to validator
		return s.sendToValidator(failedTransaction)
	}

	// Create successful transaction
	transaction := &models.Transaction{
		ID:           generateTransactionID(),
		UserID:       data.UserID,
		Type:         models.Buy,
		CurrencyPair: data.CurrencyPair,
		Amount:       data.Amount,
		ExchangeRate: quotation.BuyPrice,
		TotalValue:   totalCost,
		Status:       "completed",
		Timestamp:    data.Timestamp,
		QuotationID:  quotation.ID,
	}

	// Save transaction to MongoDB
	if err := s.currencyTransactionRepo.Create(ctx, transaction); err != nil {
		return fmt.Errorf("failed to save transaction: %v", err)
	}

	// Update user balance in PostgreSQL
	user.Balance -= totalCost
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user balance: %v", err)
	}

	// Send transaction to validation service
	return s.sendToValidator(transaction)
}

func (s *TransactionService) handleSell(ctx context.Context, data events.TransactionEventData) error {
	// Convert user ID to uint
	userID, err := strconv.ParseUint(data.UserID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	// Get user to update balance after transaction
	user, err := s.userRepo.GetByID(ctx, uint(userID))
	if err != nil {
		return fmt.Errorf("failed to find user: %v", err)
	}

	// Get latest quotation for the currency pair
	quotation, err := s.quotationRepo.GetLatestByCurrencyPair(ctx, data.CurrencyPair)
	if err != nil {
		// If no quotation found, use fallback rates
		return s.handleSellWithFallbackRate(ctx, data, user)
	}

	// Calculate total value received using real quotation
	totalValue := data.Amount * quotation.SellPrice

	// Create successful transaction
	transaction := &models.Transaction{
		ID:           generateTransactionID(),
		UserID:       data.UserID,
		Type:         models.Sell,
		CurrencyPair: data.CurrencyPair,
		Amount:       data.Amount,
		ExchangeRate: quotation.SellPrice,
		TotalValue:   totalValue,
		Status:       "completed",
		Timestamp:    data.Timestamp,
		QuotationID:  quotation.ID,
	}

	// Save transaction to MongoDB
	if err := s.currencyTransactionRepo.Create(ctx, transaction); err != nil {
		return fmt.Errorf("failed to save transaction: %v", err)
	}

	// Update user balance (add money from selling) in PostgreSQL
	user.Balance += totalValue
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user balance: %v", err)
	}

	// Send transaction to validation service
	return s.sendToValidator(transaction)
}

// Fallback methods for when no quotation is available
func (s *TransactionService) handleBuyWithFallbackRate(ctx context.Context, data events.TransactionEventData, user *models.User) error {
	// Fallback rate: 1 USD = 5.5 BRL
	fallbackRate := 5.5
	totalCost := data.Amount * fallbackRate

	// Check if user has sufficient balance
	if user.Balance < totalCost {
		// Create and save failed transaction
		failedTransaction := &models.Transaction{
			ID:           generateTransactionID(),
			UserID:       data.UserID,
			Type:         models.Buy,
			CurrencyPair: data.CurrencyPair,
			Amount:       data.Amount,
			ExchangeRate: fallbackRate,
			TotalValue:   totalCost,
			Status:       "failed: insufficient_balance",
			Timestamp:    data.Timestamp,
			QuotationID:  "fallback",
		}

		// Save failed transaction to MongoDB
		if err := s.currencyTransactionRepo.Create(ctx, failedTransaction); err != nil {
			return fmt.Errorf("failed to save failed transaction: %v", err)
		}

		// Send to validator
		return s.sendToValidator(failedTransaction)
	}

	// Create successful transaction with fallback rate
	transaction := &models.Transaction{
		ID:           generateTransactionID(),
		UserID:       data.UserID,
		Type:         models.Buy,
		CurrencyPair: data.CurrencyPair,
		Amount:       data.Amount,
		ExchangeRate: fallbackRate,
		TotalValue:   totalCost,
		Status:       "completed_fallback_rate",
		Timestamp:    data.Timestamp,
		QuotationID:  "fallback",
	}

	// Save transaction to MongoDB
	if err := s.currencyTransactionRepo.Create(ctx, transaction); err != nil {
		return fmt.Errorf("failed to save transaction: %v", err)
	}

	// Update user balance in PostgreSQL
	user.Balance -= totalCost
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user balance: %v", err)
	}

	// Send transaction to validation service
	return s.sendToValidator(transaction)
}

func (s *TransactionService) handleSellWithFallbackRate(ctx context.Context, data events.TransactionEventData, user *models.User) error {
	// Fallback rate: 1 USD = 5.3 BRL (with spread)
	fallbackRate := 5.3
	totalValue := data.Amount * fallbackRate

	// Create successful transaction with fallback rate
	transaction := &models.Transaction{
		ID:           generateTransactionID(),
		UserID:       data.UserID,
		Type:         models.Sell,
		CurrencyPair: data.CurrencyPair,
		Amount:       data.Amount,
		ExchangeRate: fallbackRate,
		TotalValue:   totalValue,
		Status:       "completed_fallback_rate",
		Timestamp:    data.Timestamp,
		QuotationID:  "fallback",
	}

	// Save transaction to MongoDB
	if err := s.currencyTransactionRepo.Create(ctx, transaction); err != nil {
		return fmt.Errorf("failed to save transaction: %v", err)
	}

	// Update user balance (add money from selling) in PostgreSQL
	user.Balance += totalValue
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user balance: %v", err)
	}

	// Send transaction to validation service
	return s.sendToValidator(transaction)
}

func (s *TransactionService) sendToValidator(transaction *models.Transaction) error {
	// Send transaction to validation service for logging
	return s.rabbitmq.PublishMessage("transactions-validator", transaction)
}

func generateTransactionID() string {
	// Generate a simple transaction ID based on timestamp
	return fmt.Sprintf("TXN_%d", time.Now().UnixNano())
}

// GetUserTransactions retrieves currency transactions for a specific user
func (s *TransactionService) GetUserTransactions(ctx context.Context, userID string, limit int) ([]*models.Transaction, error) {
	return s.currencyTransactionRepo.GetByUserID(ctx, userID, limit)
}
