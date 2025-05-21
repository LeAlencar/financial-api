package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leandroalencar/banco-dados/services/s1-generator/internal/infra/database/db"
	"github.com/leandroalencar/banco-dados/shared/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
		q:    db.New(pool),
	}
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	dbUser, err := r.q.GetUser(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	// Convert Numeric to float64
	var balance float64
	dbUser.Balance.Scan(&balance)

	return &models.User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		Balance:   balance,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	dbUser, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Convert Numeric to float64
	var balance float64
	dbUser.Balance.Scan(&balance)

	return &models.User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		Balance:   balance,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}, nil
}
