package repositories

import (
	"context"
	"errors"

	"github.com/leandroalencar/banco-dados/shared/models"
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
func (r *AccountRepository) Create(ctx context.Context, account *models.User) error {
	return r.db.WithContext(ctx).Create(account).Error
}

// GetByID retrieves an account by ID
func (r *AccountRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var account models.User
	if err := r.db.WithContext(ctx).First(&account, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return &account, nil
}

// GetAllByUserID retrieves all accounts for a user
func (r *AccountRepository) GetAllByUserID(ctx context.Context, userID string) ([]models.User, error) {
	var accounts []models.User
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

// Update updates an existing account
func (r *AccountRepository) Update(ctx context.Context, account *models.User) error {
	return r.db.WithContext(ctx).Save(account).Error
}

// Delete removes an account from the database
func (r *AccountRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

// UpdateBalance updates the balance of an account
func (r *AccountRepository) UpdateBalance(ctx context.Context, id string, amount float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var account models.User
		if err := tx.WithContext(ctx).First(&account, "id = ?", id).Error; err != nil {
			return err
		}

		account.Balance += amount

		return tx.WithContext(ctx).Save(&account).Error
	})
}
