package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/leandroalencar/banco-dados/shared/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// QuotationRepository handles MongoDB operations for quotations
type QuotationRepository struct {
	collection *mongo.Collection
}

// NewQuotationRepository creates a new quotation repository
func NewQuotationRepository(db *mongo.Database) *QuotationRepository {
	return &QuotationRepository{
		collection: db.Collection("quotations"),
	}
}

// Create saves a new quotation to MongoDB
func (r *QuotationRepository) Create(ctx context.Context, quotation *models.Quotation) error {
	// Set timestamps
	//quotation.CreatedAt = time.Now()
	//quotation.UpdatedAt = time.Now()

	// Insert the document
	_, err := r.collection.InsertOne(ctx, quotation)
	if err != nil {
		return fmt.Errorf("failed to insert quotation: %v", err)
	}

	return nil
}

// GetByID retrieves a quotation by its ID
func (r *QuotationRepository) GetByID(ctx context.Context, id string) (*models.Quotation, error) {
	var quotation models.Quotation
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&quotation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("quotation not found")
		}
		return nil, fmt.Errorf("failed to get quotation: %v", err)
	}
	return &quotation, nil
}

// GetByDateRange retrieves quotations within a date range
func (r *QuotationRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*models.Quotation, error) {
	filter := bson.M{
		"created_at": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query quotations: %v", err)
	}
	defer cursor.Close(ctx)

	var quotations []*models.Quotation
	if err := cursor.All(ctx, &quotations); err != nil {
		return nil, fmt.Errorf("failed to decode quotations: %v", err)
	}

	return quotations, nil
}
