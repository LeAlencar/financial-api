package repositories

import (
	"context"
	"fmt"

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
