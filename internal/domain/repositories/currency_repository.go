package repositories

import (
	"context"

	"github.com/lealencar/financial-api/internal/domain/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type CurrencyRepository struct {
	collection *mongo.Collection
}

func NewCurrencyRepository(db *mongo.Database) *CurrencyRepository {
	return &CurrencyRepository{
		collection: db.Collection("exchange_rates"),
	}
}

func (r *CurrencyRepository) Insert(ctx context.Context, currency *models.Currency) error {
	_, err := r.collection.InsertOne(ctx, currency)
	return err
}
