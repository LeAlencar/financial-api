package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leandroalencar/banco-dados/services/s2-processor/internal/infra/database/db"
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

// Create adds a new user to the database
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	// Always set initial balance to 0
	balance := pgtype.Numeric{}
	balance.Scan(0.00) // Set balance to 0

	// Convert time.Time to Timestamptz
	createdAt := pgtype.Timestamptz{
		Time:  user.CreatedAt,
		Valid: true,
	}

	updatedAt := pgtype.Timestamptz{
		Time:  user.UpdatedAt,
		Valid: true,
	}

	dbUser, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Balance:   balance,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return err
	}

	// Update the user's ID with the one from the database
	user.ID = dbUser.ID
	return nil
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

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	balance := pgtype.Numeric{}
	balance.Scan(user.Balance)

	updatedAt := pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	}

	_, err := r.q.UpdateUser(ctx, db.UpdateUserParams{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Balance:   balance,
		UpdatedAt: updatedAt,
	})
	return err
}

// Delete removes a user from the database
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteUser(ctx, int32(id))
}
