package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

// TransactionRepository handles database operations for transactions
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create adds a new transaction to the database
func (r *TransactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(ctx context.Context, id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.WithContext(ctx).First(&transaction, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &transaction, nil
}

// GetAllByAccountID retrieves all transactions for an account
func (r *TransactionRepository) GetAllByAccountID(ctx context.Context, accountID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.WithContext(ctx).Where("account_id = ?", accountID).Order("date DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// GetAllByUserID retrieves all transactions for a user with optional filters
func (r *TransactionRepository) GetAllByUserID(ctx context.Context, userID uint, filters map[string]interface{}) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := r.db.WithContext(ctx).
		Table("transactions").
		Joins("JOIN accounts ON transactions.account_id = accounts.id").
		Where("accounts.user_id = ?", userID)

	// Apply filters if provided
	if categoryID, ok := filters["category_id"].(uint); ok && categoryID > 0 {
		query = query.Where("transactions.category_id = ?", categoryID)
	}

	if transactionType, ok := filters["type"].(string); ok && transactionType != "" {
		query = query.Where("transactions.type = ?", transactionType)
	}

	if startDate, ok := filters["start_date"].(time.Time); ok {
		query = query.Where("transactions.date >= ?", startDate)
	}

	if endDate, ok := filters["end_date"].(time.Time); ok {
		query = query.Where("transactions.date <= ?", endDate)
	}

	if err := query.Order("transactions.date DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetByDateRange retrieves transactions within a date range for a user
func (r *TransactionRepository) GetByDateRange(ctx context.Context, userID uint, startDate, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction

	if err := r.db.WithContext(ctx).
		Table("transactions").
		Joins("JOIN accounts ON transactions.account_id = accounts.id").
		Where("accounts.user_id = ? AND transactions.date BETWEEN ? AND ?", userID, startDate, endDate).
		Order("transactions.date DESC").
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

// Update updates an existing transaction
func (r *TransactionRepository) Update(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

// Delete removes a transaction from the database
func (r *TransactionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Transaction{}, id).Error
}

// TransactionRepository handles MongoDB operations for transactions
type TransactionRepositoryMongo struct {
	collection *mongo.Collection
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepositoryMongo(db *mongo.Database) *TransactionRepositoryMongo {
	return &TransactionRepositoryMongo{
		collection: db.Collection("currency_transactions"),
	}
}

// GetByUserID retrieves transactions for a specific user
func (r *TransactionRepositoryMongo) GetByUserID(ctx context.Context, userID string, limit int) ([]*models.Transaction, error) {
	filter := bson.M{"user_id": userID}

	// Sort by timestamp descending to get the latest transactions
	opts := options.Find().SetSort(bson.D{{"timestamp", -1}})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %v", err)
	}
	defer cursor.Close(ctx)

	var transactions []*models.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, fmt.Errorf("failed to decode transactions: %v", err)
	}

	return transactions, nil
}
