package services

import (
	"context"
	stderrors "errors"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/atreya2011/expense-manager/internal/clock"
	"github.com/atreya2011/expense-manager/internal/errors"
	"github.com/atreya2011/expense-manager/internal/log"
	"github.com/atreya2011/expense-manager/internal/repo"
	db "github.com/atreya2011/expense-manager/internal/repo/gen"
	expensesv1 "github.com/atreya2011/expense-manager/internal/rpc/gen/expenses/v1"
	"github.com/atreya2011/expense-manager/internal/rpc/gen/expenses/v1/expensesv1connect"
)

// UserService implements the UserService interface defined in the proto
type UserService struct {
	expensesv1connect.UnimplementedUserServiceHandler
	repo   *repo.UserRepo
	clock  clock.Clock
	logger *slog.Logger
}

// NewUserService creates a new UserService
func NewUserService(repo *repo.UserRepo, clock clock.Clock, logger *slog.Logger) *UserService {
	return &UserService{
		repo:   repo,
		clock:  clock,
		logger: logger,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *connect.Request[expensesv1.CreateUserRequest]) (*connect.Response[expensesv1.CreateUserResponse], error) {
	// Log method entry with context
	log.InfoContext(ctx, s.logger, "Creating user", "name", req.Msg.Name, "email", req.Msg.Email)

	// Validate input
	if req.Msg.Name == "" || req.Msg.Email == "" {
		log.ErrorContext(ctx, s.logger, "Invalid input for CreateUser", "error", "name and email are required")
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: name and email are required", errors.ErrInvalidInput))
	}

	// Begin transaction
	tx, err := s.repo.GetDB().BeginTxx(ctx, nil)
	if err != nil {
		log.ErrorContext(ctx, s.logger, "Failed to begin transaction", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: failed to begin transaction: %v", errors.ErrInternal, err))
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			// Cannot return error from defer, just log it
			log.ErrorContext(ctx, s.logger, "Failed to rollback transaction", "error", err)
		}
	}() // Rollback if any error occurs

	// Create user in database within the transaction
	user, err := s.repo.CreateUser(ctx, tx, db.CreateUserParams{
		Name:  req.Msg.Name,
		Email: req.Msg.Email,
	})
	if err != nil {
		if stderrors.Is(err, errors.ErrDuplicate) {
			log.ErrorContext(ctx, s.logger, "User already exists", "email", req.Msg.Email)
			return nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("%w: user with email %s already exists", errors.ErrDuplicate, req.Msg.Email))
		}
		log.ErrorContext(ctx, s.logger, "Failed to create user", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.ErrorContext(ctx, s.logger, "Failed to commit transaction", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: failed to commit transaction: %v", errors.ErrInternal, err))
	}

	// Log success
	log.InfoContext(ctx, s.logger, "User created successfully", "id", user.ID)

	// Prepare response
	return connect.NewResponse(&expensesv1.CreateUserResponse{
		User: toProtoUser(user),
	}), nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, req *connect.Request[expensesv1.GetUserRequest]) (*connect.Response[expensesv1.GetUserResponse], error) {
	// Log method entry
	log.InfoContext(ctx, s.logger, "Getting user", "id", req.Msg.Id)

	// Validate input
	if req.Msg.Id == "" {
		log.ErrorContext(ctx, s.logger, "Invalid input for GetUser", "error", "id is required")
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: id is required", errors.ErrInvalidInput))
	}

	// Get user from database (read operations can use the main DB connection)
	user, err := s.repo.GetUser(ctx, s.repo.GetDB(), req.Msg.Id)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			log.ErrorContext(ctx, s.logger, "User not found", "id", req.Msg.Id)
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("%w: user with id %s not found", errors.ErrNotFound, req.Msg.Id))
		}
		log.ErrorContext(ctx, s.logger, "Failed to get user", "id", req.Msg.Id, "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	log.InfoContext(ctx, s.logger, "User retrieved successfully", "id", user.ID)

	// Prepare response
	return connect.NewResponse(&expensesv1.GetUserResponse{
		User: toProtoUser(user),
	}), nil
}

// ListUsers retrieves a paginated list of users
func (s *UserService) ListUsers(ctx context.Context, req *connect.Request[expensesv1.ListUsersRequest]) (*connect.Response[expensesv1.ListUsersResponse], error) {
	// Log method entry
	log.InfoContext(ctx, s.logger, "Listing users")

	// Parse pagination parameters
	limit := int32(50) // default limit
	offset := int32(0) // default offset
	if req.Msg.Pagination != nil {
		if req.Msg.Pagination.PageSize > 0 {
			limit = req.Msg.Pagination.PageSize
		}
		// Extract offset from page token if provided
		if req.Msg.Pagination.PageToken != "" {
			// In a real implementation, we would decode the page token to get the offset
			// For simplicity, we'll just parse it as an offset
			if _, err := fmt.Sscanf(req.Msg.Pagination.PageToken, "%d", &offset); err != nil {
				log.ErrorContext(ctx, s.logger, "Invalid page token", "token", req.Msg.Pagination.PageToken, "error", err)
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: invalid page token", errors.ErrInvalidInput))
			}
		}
	}

	log.InfoContext(ctx, s.logger, "Pagination parameters", "limit", limit, "offset", offset)

	// Get users from database (read operations can use the main DB connection)
	users, err := s.repo.ListUsers(ctx, s.repo.GetDB(), db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.ErrorContext(ctx, s.logger, "Failed to list users", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Prepare response
	protoUsers := make([]*expensesv1.User, len(users))
	for i, user := range users {
		protoUsers[i] = toProtoUser(user)
	}

	// Prepare pagination response
	nextPageToken := ""
	if len(users) == int(limit) {
		nextPageToken = fmt.Sprintf("%d", offset+limit)
	}

	log.InfoContext(ctx, s.logger, "Users retrieved successfully", "count", len(users))

	return connect.NewResponse(&expensesv1.ListUsersResponse{
		Users: protoUsers,
		PaginationResponse: &expensesv1.PaginationResponse{
			NextPageToken: nextPageToken,
			TotalCount:    int32(len(users)), // In a real implementation, we would get the total count from the database
		},
	}), nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, req *connect.Request[expensesv1.UpdateUserRequest]) (*connect.Response[expensesv1.UpdateUserResponse], error) {
	// Log method entry
	log.InfoContext(ctx, s.logger, "Updating user", "id", req.Msg.Id, "name", req.Msg.Name, "email", req.Msg.Email)

	// Validate input
	if req.Msg.Id == "" {
		log.ErrorContext(ctx, s.logger, "Invalid input for UpdateUser", "error", "id is required")
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: id is required", errors.ErrInvalidInput))
	}
	if req.Msg.Name == "" || req.Msg.Email == "" {
		log.ErrorContext(ctx, s.logger, "Invalid input for UpdateUser", "error", "name and email are required")
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: name and email are required", errors.ErrInvalidInput))
	}

	// Begin transaction
	tx, err := s.repo.GetDB().BeginTxx(ctx, nil)
	if err != nil {
		log.ErrorContext(ctx, s.logger, "Failed to begin transaction", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: failed to begin transaction: %v", errors.ErrInternal, err))
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			// Cannot return error from defer, just log it
			log.ErrorContext(ctx, s.logger, "Failed to rollback transaction", "error", err)
		}
	}() // Rollback if any error occurs

	// Check if user exists within the transaction
	_, err = s.repo.GetUser(ctx, tx, req.Msg.Id)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			log.ErrorContext(ctx, s.logger, "User not found", "id", req.Msg.Id)
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("%w: user with id %s not found", errors.ErrNotFound, req.Msg.Id))
		}
		log.ErrorContext(ctx, s.logger, "Failed to check if user exists", "id", req.Msg.Id, "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Update user in database within the transaction
	user, err := s.repo.UpdateUser(ctx, tx, db.UpdateUserParams{
		ID:    req.Msg.Id,
		Name:  req.Msg.Name,
		Email: req.Msg.Email,
	})
	if err != nil {
		if stderrors.Is(err, errors.ErrDuplicate) {
			log.ErrorContext(ctx, s.logger, "User with email already exists", "email", req.Msg.Email)
			return nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("%w: user with email %s already exists", errors.ErrDuplicate, req.Msg.Email))
		}
		log.ErrorContext(ctx, s.logger, "Failed to update user", "id", req.Msg.Id, "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.ErrorContext(ctx, s.logger, "Failed to commit transaction", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: failed to commit transaction: %v", errors.ErrInternal, err))
	}

	// Log success
	log.InfoContext(ctx, s.logger, "User updated successfully", "id", user.ID)

	// Prepare response
	return connect.NewResponse(&expensesv1.UpdateUserResponse{
		User: toProtoUser(user),
	}), nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, req *connect.Request[expensesv1.DeleteUserRequest]) (*connect.Response[expensesv1.DeleteUserResponse], error) {
	// Log method entry
	log.InfoContext(ctx, s.logger, "Deleting user", "id", req.Msg.Id)

	// Validate input
	if req.Msg.Id == "" {
		log.ErrorContext(ctx, s.logger, "Invalid input for DeleteUser", "error", "id is required")
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: id is required", errors.ErrInvalidInput))
	}

	// Begin transaction
	tx, err := s.repo.GetDB().BeginTxx(ctx, nil)
	if err != nil {
		log.ErrorContext(ctx, s.logger, "Failed to begin transaction", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: failed to begin transaction: %v", errors.ErrInternal, err))
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			// Cannot return error from defer, just log it
			log.ErrorContext(ctx, s.logger, "Failed to rollback transaction", "error", err)
		}
	}() // Rollback if any error occurs

	// Check if user exists within the transaction
	_, err = s.repo.GetUser(ctx, tx, req.Msg.Id)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			log.ErrorContext(ctx, s.logger, "User not found", "id", req.Msg.Id)
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("%w: user with id %s not found", errors.ErrNotFound, req.Msg.Id))
		}
		log.ErrorContext(ctx, s.logger, "Failed to check if user exists", "id", req.Msg.Id, "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Delete user from database within the transaction
	err = s.repo.DeleteUser(ctx, tx, req.Msg.Id)
	if err != nil {
		log.ErrorContext(ctx, s.logger, "Failed to delete user", "id", req.Msg.Id, "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.ErrorContext(ctx, s.logger, "Failed to commit transaction", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: failed to commit transaction: %v", errors.ErrInternal, err))
	}

	// Log success
	log.InfoContext(ctx, s.logger, "User deleted successfully", "id", req.Msg.Id)

	// Prepare response
	return connect.NewResponse(&expensesv1.DeleteUserResponse{
		Success: true,
	}), nil
}

// toProtoUser converts a db.User to a expensesv1.User
func toProtoUser(user db.User) *expensesv1.User {
	return &expensesv1.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
