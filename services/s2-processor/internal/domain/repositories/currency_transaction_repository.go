package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/leandroalencar/banco-dados/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CurrencyTransactionRepository handles MongoDB operations for currency transactions
type CurrencyTransactionRepository struct {
	collection *mongo.Collection
}

// NewCurrencyTransactionRepository creates a new currency transaction repository
func NewCurrencyTransactionRepository(db *mongo.Database) *CurrencyTransactionRepository {
	return &CurrencyTransactionRepository{
		collection: db.Collection("currency_transactions"),
	}
}

// Create saves a new currency transaction to MongoDB
func (r *CurrencyTransactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	// Insert the document
	_, err := r.collection.InsertOne(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to insert currency transaction: %v", err)
	}

	return nil
}

// GetByID retrieves a currency transaction by its ID
func (r *CurrencyTransactionRepository) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("currency transaction not found")
		}
		return nil, fmt.Errorf("failed to get currency transaction: %v", err)
	}
	return &transaction, nil
}

// GetByUserID retrieves currency transactions for a specific user
func (r *CurrencyTransactionRepository) GetByUserID(ctx context.Context, userID string, limit int) ([]*models.Transaction, error) {
	filter := bson.M{"user_id": userID}

	// Sort by timestamp descending to get the latest transactions
	opts := options.Find().SetSort(bson.D{{"timestamp", -1}})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query currency transactions: %v", err)
	}
	defer cursor.Close(ctx)

	var transactions []*models.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, fmt.Errorf("failed to decode currency transactions: %v", err)
	}

	return transactions, nil
}

// GetByDateRange retrieves currency transactions within a date range for a user
func (r *CurrencyTransactionRepository) GetByDateRange(ctx context.Context, userID string, start, end time.Time) ([]*models.Transaction, error) {
	filter := bson.M{
		"user_id": userID,
		"timestamp": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query currency transactions by date range: %v", err)
	}
	defer cursor.Close(ctx)

	var transactions []*models.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, fmt.Errorf("failed to decode currency transactions: %v", err)
	}

	return transactions, nil
}

// GetByStatus retrieves currency transactions by status
func (r *CurrencyTransactionRepository) GetByStatus(ctx context.Context, status string, limit int) ([]*models.Transaction, error) {
	filter := bson.M{"status": status}

	opts := options.Find().SetSort(bson.D{{"timestamp", -1}})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query currency transactions by status: %v", err)
	}
	defer cursor.Close(ctx)

	var transactions []*models.Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, fmt.Errorf("failed to decode currency transactions: %v", err)
	}

	return transactions, nil
}

// Update updates an existing currency transaction
func (r *CurrencyTransactionRepository) Update(ctx context.Context, transaction *models.Transaction) error {
	filter := bson.M{"_id": transaction.ID}
	update := bson.M{"$set": transaction}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update currency transaction: %v", err)
	}

	return nil
}
