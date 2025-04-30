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

// InstrumentService implements the InstrumentService interface defined in the proto
type InstrumentService struct {
	expensesv1connect.UnimplementedInstrumentServiceHandler
	repo  *repo.InstrumentRepo
	clock clock.Clock
}

// NewInstrumentService creates a new InstrumentService
func NewInstrumentService(repo *repo.InstrumentRepo, clock clock.Clock) *InstrumentService {
	return &InstrumentService{
		repo:  repo,
		clock: clock,
	}
}

// CreateInstrument creates a new instrument
func (s *InstrumentService) CreateInstrument(ctx context.Context, req *connect.Request[expensesv1.CreateInstrumentRequest]) (*connect.Response[expensesv1.CreateInstrumentResponse], error) {
	// Validate input
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: name is required", errors.ErrInvalidInput))
	}

	// Create instrument in database
	instrument, err := s.repo.CreateInstrument(ctx, req.Msg.Name)
	if err != nil {
		if stderrors.Is(err, errors.ErrDuplicate) {
			return nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("instrument with name '%s' already exists", req.Msg.Name))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Prepare response
	return connect.NewResponse(&expensesv1.CreateInstrumentResponse{
		Instrument: toProtoInstrument(instrument),
	}), nil
}

// GetInstrument retrieves an instrument by ID
func (s *InstrumentService) GetInstrument(ctx context.Context, req *connect.Request[expensesv1.GetInstrumentRequest]) (*connect.Response[expensesv1.GetInstrumentResponse], error) {
	// Validate input
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: id is required", errors.ErrInvalidInput))
	}

	// Get instrument from database
	instrument, err := s.repo.GetInstrument(ctx, req.Msg.Id)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("%w: instrument with id %s not found", errors.ErrNotFound, req.Msg.Id))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Prepare response
	return connect.NewResponse(&expensesv1.GetInstrumentResponse{
		Instrument: toProtoInstrument(instrument),
	}), nil
}

// ListInstruments retrieves a paginated list of instruments
func (s *InstrumentService) ListInstruments(ctx context.Context, req *connect.Request[expensesv1.ListInstrumentsRequest]) (*connect.Response[expensesv1.ListInstrumentsResponse], error) {
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

	// Get instruments from database
	instruments, err := s.repo.ListInstruments(ctx, int64(limit), int64(offset))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Prepare response
	protoInstruments := make([]*expensesv1.Instrument, len(instruments))
	for i, instrument := range instruments {
		protoInstruments[i] = toProtoInstrument(instrument)
	}

	// Prepare pagination response
	nextPageToken := ""
	if len(instruments) == int(limit) {
		nextPageToken = fmt.Sprintf("%d", offset+limit)
	}

	return connect.NewResponse(&expensesv1.ListInstrumentsResponse{
		Instruments: protoInstruments,
		PaginationResponse: &expensesv1.PaginationResponse{
			NextPageToken: nextPageToken,
			TotalCount:    int32(len(instruments)), // In a real implementation, we would get the total count from the database
		},
	}), nil
}

// UpdateInstrument updates an instrument
func (s *InstrumentService) UpdateInstrument(ctx context.Context, req *connect.Request[expensesv1.UpdateInstrumentRequest]) (*connect.Response[expensesv1.UpdateInstrumentResponse], error) {
	// Validate input
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: id is required", errors.ErrInvalidInput))
	}
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: name is required", errors.ErrInvalidInput))
	}

	// Check if instrument exists
	_, err := s.repo.GetInstrument(ctx, req.Msg.Id)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("%w: instrument with id %s not found", errors.ErrNotFound, req.Msg.Id))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Update instrument in database
	instrument, err := s.repo.UpdateInstrument(ctx, req.Msg.Id, req.Msg.Name)
	if err != nil {
		if stderrors.Is(err, errors.ErrDuplicate) {
			return nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("instrument with name '%s' already exists", req.Msg.Name))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Prepare response
	return connect.NewResponse(&expensesv1.UpdateInstrumentResponse{
		Instrument: toProtoInstrument(instrument),
	}), nil
}

// DeleteInstrument deletes an instrument by ID
func (s *InstrumentService) DeleteInstrument(ctx context.Context, req *connect.Request[expensesv1.DeleteInstrumentRequest]) (*connect.Response[expensesv1.DeleteInstrumentResponse], error) {
	// Validate input
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%w: id is required", errors.ErrInvalidInput))
	}

	// Check if instrument exists
	_, err := s.repo.GetInstrument(ctx, req.Msg.Id)
	if err != nil {
		if stderrors.Is(err, errors.ErrNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("%w: instrument with id %s not found", errors.ErrNotFound, req.Msg.Id))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Delete instrument from database
	err = s.repo.DeleteInstrument(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("%w: %v", errors.ErrInternal, err))
	}

	// Prepare response
	return connect.NewResponse(&expensesv1.DeleteInstrumentResponse{
		Success: true,
	}), nil
}

// toProtoInstrument converts a db.Instrument to a expensesv1.Instrument
func toProtoInstrument(instrument db.Instrument) *expensesv1.Instrument {
	return &expensesv1.Instrument{
		Id:        instrument.ID,
		Name:      instrument.Name,
		CreatedAt: timestamppb.New(instrument.CreatedAt),
		UpdatedAt: timestamppb.New(instrument.UpdatedAt),
	}
}
