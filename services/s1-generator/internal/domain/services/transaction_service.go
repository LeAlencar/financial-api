package services

import (
	"context"

	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/domain/repositories"
	"github.com/leandroalencar/banco-dados/shared/models"
)

// TransactionService handles business logic related to currency transactions
type TransactionService struct {
	currencyTransactionRepo *repositories.CurrencyTransactionRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(currencyTransactionRepo *repositories.CurrencyTransactionRepository) *TransactionService {
	return &TransactionService{
		currencyTransactionRepo: currencyTransactionRepo,
	}
}

// GetUserTransactions retrieves currency transactions for a specific user
func (s *TransactionService) GetUserTransactions(ctx context.Context, userID string, limit int) ([]*models.Transaction, error) {
	return s.currencyTransactionRepo.GetByUserID(ctx, userID, limit)
}
