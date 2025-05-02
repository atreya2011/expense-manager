package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atreya2011/expense-manager/internal/errors"
	db "github.com/atreya2011/expense-manager/internal/repo/gen"
	"github.com/jmoiron/sqlx"
)

// UserRepo provides direct access to user-related database operations
type UserRepo struct {
	db *sqlx.DB // Store the underlying DB pool
}

// NewUserRepo creates a new UserRepo
func NewUserRepo(dbConn *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: dbConn,
	}
}

// GetDB returns the underlying database connection pool
func (r *UserRepo) GetDB() *sqlx.DB {
	return r.db
}

// CreateUser creates a new user within the provided DBTX
func (r *UserRepo) CreateUser(ctx context.Context, dbtx db.DBTX, arg db.CreateUserParams) (db.User, error) {
	queries := db.New(dbtx)
	user, err := queries.CreateUser(ctx, arg)
	if err != nil {
		return db.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// GetUser retrieves a user by ID within the provided DBTX
func (r *UserRepo) GetUser(ctx context.Context, dbtx db.DBTX, id string) (db.User, error) {
	queries := db.New(dbtx)
	user, err := queries.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, fmt.Errorf("user not found: %w", errors.ErrNotFound)
		}
		return db.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// ListUsers retrieves a paginated list of users within the provided DBTX
func (r *UserRepo) ListUsers(ctx context.Context, dbtx db.DBTX, arg db.ListUsersParams) ([]db.User, error) {
	queries := db.New(dbtx)
	users, err := queries.ListUsers(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// UpdateUser updates a user within the provided DBTX
func (r *UserRepo) UpdateUser(ctx context.Context, dbtx db.DBTX, arg db.UpdateUserParams) (db.User, error) {
	queries := db.New(dbtx)
	user, err := queries.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, fmt.Errorf("user not found: %w", errors.ErrNotFound)
		}
		return db.User{}, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

// DeleteUser deletes a user by ID within the provided DBTX
func (r *UserRepo) DeleteUser(ctx context.Context, dbtx db.DBTX, id string) error {
	queries := db.New(dbtx)
	err := queries.DeleteUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found: %w", errors.ErrNotFound)
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
