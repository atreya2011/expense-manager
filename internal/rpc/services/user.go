package services

import (
	"context"
	stderrors "errors"
	"fmt"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/atreya2011/expense-manager/internal/clock"
	"github.com/atreya2011/expense-manager/internal/errors"
	"github.com/atreya2011/expense-manager/internal/repo"
	db "github.com/atreya2011/expense-manager/internal/repo/gen"
	expensesv1 "github.com/atreya2011/expense-manager/internal/rpc/gen/expenses/v1"
	"github.com/atreya2011/expense-manager/internal/rpc/gen/expenses/v1/expensesv1connect"
)

// UserService implements the UserService interface defined in the proto
type UserService struct {
	expensesv1connect.UnimplementedUserServiceHandler
	repo  *repo.UserRepo
	clock clock.Clock
}

// NewUserService creates a new UserService
func NewUserService(repo *repo.UserRepo, clock clock.Clock) *UserService {
	return &UserService{
		repo:  repo,
		clock: clock,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *connect.Request[expensesv1.CreateUserRequest]) (*connect.Response[expensesv1.CreateUserResponse], error) {
	// Validate input
	if req.Msg.Name == "" || req.Msg.Email == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: name and email are required", errors.ErrInvalidInput))
	}

	// Create user in database
	user, err := s.repo.CreateUser(ctx, db.CreateUserParams{
		Name:  req.Msg.Name,
		Email: req.Msg.Email,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Prepare response
	return connect.NewResponse(&expensesv1.CreateUserResponse{
		User: toProtoUser(user),
	}), nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, req *connect.Request[expensesv1.GetUserRequest]) (*connect.Response[expensesv1.GetUserResponse], error) {
	// Validate input
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: id is required", errors.ErrInvalidInput))
	}

	// Get user from database
	user, err := s.repo.GetUser(ctx, req.Msg.Id)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("%w: user with id %s not found", errors.ErrNotFound, req.Msg.Id))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Prepare response
	return connect.NewResponse(&expensesv1.GetUserResponse{
		User: toProtoUser(user),
	}), nil
}

// ListUsers retrieves a paginated list of users
func (s *UserService) ListUsers(ctx context.Context, req *connect.Request[expensesv1.ListUsersRequest]) (*connect.Response[expensesv1.ListUsersResponse], error) {
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
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: invalid page token", errors.ErrInvalidInput))
			}
		}
	}

	// Get users from database
	users, err := s.repo.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
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
	// Validate input
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: id is required", errors.ErrInvalidInput))
	}
	if req.Msg.Name == "" || req.Msg.Email == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: name and email are required", errors.ErrInvalidInput))
	}

	// Check if user exists
	_, err := s.repo.GetUser(ctx, req.Msg.Id)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("%w: user with id %s not found", errors.ErrNotFound, req.Msg.Id))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Update user in database
	user, err := s.repo.UpdateUser(ctx, db.UpdateUserParams{
		ID:    req.Msg.Id,
		Name:  req.Msg.Name,
		Email: req.Msg.Email,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Prepare response
	return connect.NewResponse(&expensesv1.UpdateUserResponse{
		User: toProtoUser(user),
	}), nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, req *connect.Request[expensesv1.DeleteUserRequest]) (*connect.Response[expensesv1.DeleteUserResponse], error) {
	// Validate input
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: id is required", errors.ErrInvalidInput))
	}

	// Check if user exists
	_, err := s.repo.GetUser(ctx, req.Msg.Id)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("%w: user with id %s not found", errors.ErrNotFound, req.Msg.Id))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Delete user from database
	err = s.repo.DeleteUser(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

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
