package repositories

import (
	"context"
	"math"
	"math/big"
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
	// Set initial balance to 1000 BRL for testing purposes
	balance := pgtype.Numeric{
		Int:   big.NewInt(100000), // 1000.00 represented as 100000 (with 2 decimal places)
		Exp:   -2,                 // 2 decimal places
		Valid: true,
	}

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
	// Set the balance in the user object too
	user.Balance = 1000.00
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	dbUser, err := r.q.GetUser(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	// Convert Numeric to float64 using the simplest approach
	var balance float64
	if dbUser.Balance.Valid && dbUser.Balance.Int != nil {
		// Try using Float64 method if available
		f, _ := new(big.Float).SetInt(dbUser.Balance.Int).Float64()
		// Apply the exponent
		balance = f * math.Pow(10, float64(dbUser.Balance.Exp))
	}

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

	// Convert Numeric to float64 using the simplest approach
	var balance float64
	if dbUser.Balance.Valid && dbUser.Balance.Int != nil {
		// Try using Float64 method if available
		f, _ := new(big.Float).SetInt(dbUser.Balance.Int).Float64()
		// Apply the exponent
		balance = f * math.Pow(10, float64(dbUser.Balance.Exp))
	}

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
	// Convert float64 to pgtype.Numeric properly
	balanceInt := int64(user.Balance * 100) // Convert to cents (2 decimal places)
	balance := pgtype.Numeric{
		Int:   big.NewInt(balanceInt),
		Exp:   -2, // 2 decimal places
		Valid: true,
	}

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
