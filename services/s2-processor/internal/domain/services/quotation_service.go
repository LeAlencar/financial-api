package services

import (
	"context"
	"fmt"

	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/domain/repositories"
	"github.com/leandroalencar/banco-dados/shared/models"
)

// QuotationService handles business logic related to quotations
type QuotationService struct {
	quotationRepo *repositories.QuotationRepository
}

// NewQuotationService creates a new quotation service
func NewQuotationService(quotationRepo *repositories.QuotationRepository) *QuotationService {
	return &QuotationService{
		quotationRepo: quotationRepo,
	}
}

// SaveQuotation saves a new quotation to the database
func (s *QuotationService) SaveQuotation(ctx context.Context, quotation *models.Quotation) error {
	if err := s.validateQuotation(quotation); err != nil {
		return fmt.Errorf("invalid quotation: %v", err)
	}

	if err := s.quotationRepo.Create(ctx, quotation); err != nil {
		return fmt.Errorf("failed to save quotation: %v", err)
	}

	return nil
}

// validateQuotation performs basic validation on a quotation
func (s *QuotationService) validateQuotation(quotation *models.Quotation) error {
	if quotation == nil {
		return fmt.Errorf("quotation cannot be nil")
	}

	if quotation.CurrencyPair == "" {
		return fmt.Errorf("currency pair is required")
	}

	if quotation.BuyPrice <= 0 {
		return fmt.Errorf("buy price must be greater than 0")
	}

	if quotation.SellPrice <= 0 {
		return fmt.Errorf("sell price must be greater than 0")
	}

	return nil
}
