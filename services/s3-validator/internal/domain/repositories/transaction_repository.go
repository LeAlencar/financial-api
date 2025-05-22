package repositories

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/leandroalencar/banco-dados/services/s3-validator/internal/domain/models"
)

type TransactionRepository struct {
	session *gocql.Session
}

func NewTransactionRepository(session *gocql.Session) *TransactionRepository {
	return &TransactionRepository{
		session: session,
	}
}

func (r *TransactionRepository) Save(ctx context.Context, transaction *models.Transaction) error {
	query := `INSERT INTO transactions (id, user_id, type, currency, amount, status, timestamp) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	return r.session.Query(query,
		transaction.ID,
		transaction.UserID,
		transaction.Type,
		transaction.Currency,
		transaction.Amount,
		transaction.Status,
		transaction.Timestamp,
	).WithContext(ctx).Exec()
}

func (r *TransactionRepository) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	var transaction models.Transaction
	query := `SELECT id, user_id, type, currency, amount, status, timestamp 
			  FROM transactions WHERE id = ?`

	err := r.session.Query(query, id).WithContext(ctx).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Type,
		&transaction.Currency,
		&transaction.Amount,
		&transaction.Status,
		&transaction.Timestamp,
	)

	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, fmt.Errorf("transaction not found: %s", id)
		}
		return nil, err
	}

	return &transaction, nil
}

func (r *TransactionRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	query := `SELECT id, user_id, type, currency, amount, status, timestamp 
			  FROM transactions WHERE user_id = ?`

	iter := r.session.Query(query, userID).WithContext(ctx).Iter()
	defer iter.Close()

	for {
		var transaction models.Transaction
		if !iter.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Type,
			&transaction.Currency,
			&transaction.Amount,
			&transaction.Status,
			&transaction.Timestamp,
		) {
			break
		}
		transactions = append(transactions, &transaction)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return transactions, nil
}
