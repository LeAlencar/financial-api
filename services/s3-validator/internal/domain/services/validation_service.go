package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/leandroalencar/banco-dados/services/s3-validator/internal/domain/models"
	"github.com/leandroalencar/banco-dados/services/s3-validator/internal/domain/repositories"
	"github.com/leandroalencar/banco-dados/shared/messaging/events"
)

type ValidationService struct {
	validationLogRepo *repositories.ValidationLogRepository
	transactionRepo   *repositories.TransactionRepository
}

func NewValidationService(validationLogRepo *repositories.ValidationLogRepository, transactionRepo *repositories.TransactionRepository) *ValidationService {
	return &ValidationService{
		validationLogRepo: validationLogRepo,
		transactionRepo:   transactionRepo,
	}
}

// ValidateUserEvent validates user-related events
func (s *ValidationService) ValidateUserEvent(ctx context.Context, userEvent *events.UserEvent, rawPayload string) error {
	validationLog := &models.ValidationLog{
		ID:          gocql.TimeUUID(),
		EventType:   models.EventTypeUser,
		Action:      models.ActionType(string(userEvent.Action)),
		ProcessedAt: time.Now(),
		RawPayload:  rawPayload,
		Source:      "s1-generator",
		Details:     make(map[string]string),
	}

	// Extract user ID if available
	if userEvent.Data.ID != 0 {
		validationLog.UserID = fmt.Sprintf("%d", userEvent.Data.ID)
		validationLog.EventID = fmt.Sprintf("user_%d", userEvent.Data.ID)
	} else {
		validationLog.EventID = fmt.Sprintf("user_%s", userEvent.Data.Email)
	}

	// Validate based on action type
	switch userEvent.Action {
	case events.UserActionCreate:
		if err := s.validateUserCreate(userEvent.Data, validationLog); err != nil {
			validationLog.Status = models.ValidationFailed
			validationLog.Message = err.Error()
		} else {
			validationLog.Status = models.ValidationSuccess
			validationLog.Message = "User creation event validated successfully"
		}
	case events.UserActionUpdate:
		if err := s.validateUserUpdate(userEvent.Data, validationLog); err != nil {
			validationLog.Status = models.ValidationFailed
			validationLog.Message = err.Error()
		} else {
			validationLog.Status = models.ValidationSuccess
			validationLog.Message = "User update event validated successfully"
		}
	case events.UserActionDelete:
		if err := s.validateUserDelete(userEvent.Data, validationLog); err != nil {
			validationLog.Status = models.ValidationFailed
			validationLog.Message = err.Error()
		} else {
			validationLog.Status = models.ValidationSuccess
			validationLog.Message = "User delete event validated successfully"
		}
	default:
		validationLog.Status = models.ValidationError
		validationLog.Message = fmt.Sprintf("Unknown user action: %s", userEvent.Action)
	}

	return s.validationLogRepo.SaveValidationLog(ctx, validationLog)
}

// ValidateTransactionEvent validates transaction-related events
func (s *ValidationService) ValidateTransactionEvent(ctx context.Context, transactionEvent *events.TransactionEvent, rawPayload string) error {
	validationLog := &models.ValidationLog{
		ID:          gocql.TimeUUID(),
		EventType:   models.EventTypeTransaction,
		Action:      models.ActionType(string(transactionEvent.Action)),
		UserID:      transactionEvent.Data.UserID,
		EventID:     fmt.Sprintf("transaction_%s_%s", transactionEvent.Data.UserID, transactionEvent.Data.Timestamp.Format("20060102150405")),
		ProcessedAt: time.Now(),
		RawPayload:  rawPayload,
		Source:      "s1-generator",
		Details:     make(map[string]string),
	}

	// Validate based on action type
	switch transactionEvent.Action {
	case events.TransactionActionBuy, events.TransactionActionSell:
		if err := s.validateTransaction(transactionEvent.Data, validationLog); err != nil {
			validationLog.Status = models.ValidationFailed
			validationLog.Message = err.Error()
		} else {
			validationLog.Status = models.ValidationSuccess
			validationLog.Message = fmt.Sprintf("Transaction %s event validated successfully", transactionEvent.Action)
		}
	default:
		validationLog.Status = models.ValidationError
		validationLog.Message = fmt.Sprintf("Unknown transaction action: %s", transactionEvent.Action)
	}

	// Also save to transaction repository for audit
	if validationLog.Status == models.ValidationSuccess {
		transaction := &models.Transaction{
			ID:        validationLog.EventID,
			UserID:    transactionEvent.Data.UserID,
			Type:      models.TransactionType(string(transactionEvent.Action)),
			Currency:  transactionEvent.Data.CurrencyPair,
			Amount:    transactionEvent.Data.Amount,
			Status:    "validated",
			Timestamp: transactionEvent.Data.Timestamp,
		}

		if err := s.transactionRepo.Save(ctx, transaction); err != nil {
			// Log the error but don't fail the validation
			validationLog.Details["transaction_save_error"] = err.Error()
		}
	}

	return s.validationLogRepo.SaveValidationLog(ctx, validationLog)
}

// ValidateQuotationEvent validates quotation-related events
func (s *ValidationService) ValidateQuotationEvent(ctx context.Context, quotationEvent *events.QuotationEvent, rawPayload string) error {
	validationLog := &models.ValidationLog{
		ID:          gocql.TimeUUID(),
		EventType:   models.EventTypeQuotation,
		Action:      models.ActionType(string(quotationEvent.Action)),
		EventID:     fmt.Sprintf("quotation_%s_%s", quotationEvent.Data.CurrencyPair, quotationEvent.Data.Timestamp.Format("20060102150405")),
		ProcessedAt: time.Now(),
		RawPayload:  rawPayload,
		Source:      "s1-generator",
		Details:     make(map[string]string),
	}

	if err := s.validateQuotation(quotationEvent.Data, validationLog); err != nil {
		validationLog.Status = models.ValidationFailed
		validationLog.Message = err.Error()
	} else {
		validationLog.Status = models.ValidationSuccess
		validationLog.Message = "Quotation event validated successfully"
	}

	return s.validationLogRepo.SaveValidationLog(ctx, validationLog)
}

// Private validation methods

func (s *ValidationService) validateUserCreate(data events.UserEventData, log *models.ValidationLog) error {
	if strings.TrimSpace(data.Name) == "" {
		log.Details["error_field"] = "name"
		return fmt.Errorf("user name is required")
	}

	if strings.TrimSpace(data.Email) == "" {
		log.Details["error_field"] = "email"
		return fmt.Errorf("user email is required")
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(data.Email) {
		log.Details["error_field"] = "email"
		log.Details["invalid_email"] = data.Email
		return fmt.Errorf("invalid email format: %s", data.Email)
	}

	if len(data.Password) < 6 {
		log.Details["error_field"] = "password"
		log.Details["password_length"] = fmt.Sprintf("%d", len(data.Password))
		return fmt.Errorf("password must be at least 6 characters long")
	}

	log.Details["validated_fields"] = "name,email,password"
	return nil
}

func (s *ValidationService) validateUserUpdate(data events.UserEventData, log *models.ValidationLog) error {
	if data.ID == 0 {
		log.Details["error_field"] = "id"
		return fmt.Errorf("user ID is required for update")
	}

	if strings.TrimSpace(data.Name) == "" && strings.TrimSpace(data.Email) == "" {
		log.Details["error_field"] = "name,email"
		return fmt.Errorf("at least one field (name or email) must be provided for update")
	}

	if data.Email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(data.Email) {
			log.Details["error_field"] = "email"
			log.Details["invalid_email"] = data.Email
			return fmt.Errorf("invalid email format: %s", data.Email)
		}
	}

	log.Details["user_id"] = fmt.Sprintf("%d", data.ID)
	return nil
}

func (s *ValidationService) validateUserDelete(data events.UserEventData, log *models.ValidationLog) error {
	if data.ID == 0 {
		log.Details["error_field"] = "id"
		return fmt.Errorf("user ID is required for delete")
	}

	log.Details["user_id"] = fmt.Sprintf("%d", data.ID)
	return nil
}

func (s *ValidationService) validateTransaction(data events.TransactionEventData, log *models.ValidationLog) error {
	if strings.TrimSpace(data.UserID) == "" {
		log.Details["error_field"] = "user_id"
		return fmt.Errorf("user ID is required for transaction")
	}

	if data.Amount <= 0 {
		log.Details["error_field"] = "amount"
		log.Details["invalid_amount"] = fmt.Sprintf("%.2f", data.Amount)
		return fmt.Errorf("transaction amount must be greater than 0")
	}

	if strings.TrimSpace(data.CurrencyPair) == "" {
		log.Details["error_field"] = "currency_pair"
		return fmt.Errorf("currency pair is required for transaction")
	}

	// Validate currency pair format (e.g., USD/BRL)
	currencyRegex := regexp.MustCompile(`^[A-Z]{3}/[A-Z]{3}$`)
	if !currencyRegex.MatchString(data.CurrencyPair) {
		log.Details["error_field"] = "currency_pair"
		log.Details["invalid_currency_pair"] = data.CurrencyPair
		return fmt.Errorf("invalid currency pair format: %s", data.CurrencyPair)
	}

	log.Details["user_id"] = data.UserID
	log.Details["amount"] = fmt.Sprintf("%.2f", data.Amount)
	log.Details["currency_pair"] = data.CurrencyPair
	return nil
}

func (s *ValidationService) validateQuotation(data events.QuotationEventData, log *models.ValidationLog) error {
	if strings.TrimSpace(data.CurrencyPair) == "" {
		log.Details["error_field"] = "currency_pair"
		return fmt.Errorf("currency pair is required for quotation")
	}

	if data.BuyPrice <= 0 {
		log.Details["error_field"] = "buy_price"
		log.Details["invalid_buy_price"] = fmt.Sprintf("%.4f", data.BuyPrice)
		return fmt.Errorf("buy price must be greater than 0")
	}

	if data.SellPrice <= 0 {
		log.Details["error_field"] = "sell_price"
		log.Details["invalid_sell_price"] = fmt.Sprintf("%.4f", data.SellPrice)
		return fmt.Errorf("sell price must be greater than 0")
	}

	// Validate that buy price is higher than sell price (spread logic)
	if data.BuyPrice <= data.SellPrice {
		log.Details["error_field"] = "price_spread"
		log.Details["buy_price"] = fmt.Sprintf("%.4f", data.BuyPrice)
		log.Details["sell_price"] = fmt.Sprintf("%.4f", data.SellPrice)
		return fmt.Errorf("buy price (%.4f) must be higher than sell price (%.4f)", data.BuyPrice, data.SellPrice)
	}

	log.Details["currency_pair"] = data.CurrencyPair
	log.Details["buy_price"] = fmt.Sprintf("%.4f", data.BuyPrice)
	log.Details["sell_price"] = fmt.Sprintf("%.4f", data.SellPrice)
	return nil
}

// GetValidationLogs retrieves validation logs with various filters
func (s *ValidationService) GetValidationLogs(ctx context.Context, eventType string, status string, limit int) ([]*models.ValidationLog, error) {
	if eventType != "" {
		return s.validationLogRepo.GetValidationLogsByEventType(ctx, models.EventType(eventType), limit)
	}

	if status != "" {
		return s.validationLogRepo.GetValidationLogsByStatus(ctx, models.ValidationStatus(status), limit)
	}

	// Default: get recent failed validations
	return s.validationLogRepo.GetValidationLogsByStatus(ctx, models.ValidationFailed, limit)
}

// GetUserValidationLogs retrieves validation logs for a specific user
func (s *ValidationService) GetUserValidationLogs(ctx context.Context, userID string, limit int) ([]*models.ValidationLog, error) {
	return s.validationLogRepo.GetValidationLogsByUserID(ctx, userID, limit)
}
