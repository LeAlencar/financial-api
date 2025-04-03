package repositories

import (
	"context"
	"errors"

	"github.com/lealencar/financial-api/internal/domain/models"
	"gorm.io/gorm"
)

// AccountRepository handles database operations for accounts
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create adds a new account to the database
func (r *AccountRepository) Create(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

// GetByID retrieves an account by ID
func (r *AccountRepository) GetByID(ctx context.Context, id uint) (*models.Account, error) {
	var account models.Account
	if err := r.db.WithContext(ctx).First(&account, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return &account, nil
}

// GetAllByUserID retrieves all accounts for a user
func (r *AccountRepository) GetAllByUserID(ctx context.Context, userID uint) ([]models.Account, error) {
	var accounts []models.Account
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

// Update updates an existing account
func (r *AccountRepository) Update(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

// Delete removes an account from the database
func (r *AccountRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Account{}, id).Error
}

// UpdateBalance updates the balance of an account
func (r *AccountRepository) UpdateBalance(ctx context.Context, id uint, amount float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var account models.Account
		if err := tx.WithContext(ctx).First(&account, id).Error; err != nil {
			return err
		}

		account.Balance += amount

		return tx.WithContext(ctx).Save(&account).Error
	})
}
