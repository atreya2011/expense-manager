package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	db "github.com/atreya2011/expense-manager/internal/repo/gen"
)

var (
	// ErrNotFound is returned when a requested resource is not found
	ErrNotFound = errors.New("not found")
)

// UserRepo provides direct access to user-related database operations
type UserRepo struct {
	q *db.Queries
}

// NewUserRepo creates a new UserRepo
func NewUserRepo(dbConn *sql.DB) *UserRepo {
	return &UserRepo{
		q: db.New(dbConn),
	}
}

// CreateUser creates a new user
func (r *UserRepo) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	user, err := r.q.CreateUser(ctx, arg)
	if err != nil {
		return db.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// GetUser retrieves a user by ID
func (r *UserRepo) GetUser(ctx context.Context, id string) (db.User, error) {
	user, err := r.q.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, fmt.Errorf("user not found: %w", ErrNotFound)
		}
		return db.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// ListUsers retrieves a paginated list of users
func (r *UserRepo) ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	users, err := r.q.ListUsers(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// UpdateUser updates a user
func (r *UserRepo) UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error) {
	user, err := r.q.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.User{}, fmt.Errorf("user not found: %w", ErrNotFound)
		}
		return db.User{}, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

// DeleteUser deletes a user by ID
func (r *UserRepo) DeleteUser(ctx context.Context, id string) error {
	err := r.q.DeleteUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found: %w", ErrNotFound)
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
