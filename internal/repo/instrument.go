package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/atreya2011/expense-manager/internal/errors"
	db "github.com/atreya2011/expense-manager/internal/repo/gen"
)

// InstrumentRepo provides direct access to instrument-related database operations
type InstrumentRepo struct {
	q *db.Queries
}

// NewInstrumentRepo creates a new InstrumentRepo
func NewInstrumentRepo(dbConn *sql.DB) *InstrumentRepo {
	return &InstrumentRepo{
		q: db.New(dbConn),
	}
}

// CreateInstrument creates a new instrument
func (r *InstrumentRepo) CreateInstrument(ctx context.Context, name string) (db.Instrument, error) {
	instrument, err := r.q.CreateInstrument(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return db.Instrument{}, fmt.Errorf("instrument with this name already exists: %w", errors.ErrDuplicate)
		}
		return db.Instrument{}, fmt.Errorf("failed to create instrument: %w", err)
	}
	return instrument, nil
}

// GetInstrument retrieves a instrument by ID
func (r *InstrumentRepo) GetInstrument(ctx context.Context, id string) (db.Instrument, error) {
	instrument, err := r.q.GetInstrument(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Instrument{}, fmt.Errorf("instrument not found: %w", errors.ErrNotFound)
		}
		return db.Instrument{}, fmt.Errorf("failed to get instrument: %w", err)
	}
	return instrument, nil
}

// ListInstruments retrieves a paginated list of instruments
func (r *InstrumentRepo) ListInstruments(ctx context.Context, limit, offset int64) ([]db.Instrument, error) {
	instruments, err := r.q.ListInstruments(ctx, db.ListInstrumentsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list instruments: %w", err)
	}
	return instruments, nil
}

// UpdateInstrument updates an instrument
func (r *InstrumentRepo) UpdateInstrument(ctx context.Context, id, name string) (db.Instrument, error) {
	instrument, err := r.q.UpdateInstrument(ctx, db.UpdateInstrumentParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Instrument{}, fmt.Errorf("instrument not found: %w", errors.ErrNotFound)
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return db.Instrument{}, fmt.Errorf("instrument with this name already exists: %w", errors.ErrDuplicate)
		}
		return db.Instrument{}, fmt.Errorf("failed to update instrument: %w", err)
	}
	return instrument, nil
}

// DeleteInstrument deletes an instrument by ID
func (r *InstrumentRepo) DeleteInstrument(ctx context.Context, id string) error {
	err := r.q.DeleteInstrument(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("instrument not found: %w", errors.ErrNotFound)
		}
		return fmt.Errorf("failed to delete instrument: %w", err)
	}
	return nil
}
