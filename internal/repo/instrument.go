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
	db *sql.DB // Store the underlying DB pool
}

// NewInstrumentRepo creates a new InstrumentRepo
func NewInstrumentRepo(dbConn *sql.DB) *InstrumentRepo {
	return &InstrumentRepo{
		db: dbConn,
	}
}

// GetDB returns the underlying database connection pool
func (r *InstrumentRepo) GetDB() *sql.DB {
	return r.db
}

// CreateInstrument creates a new instrument within the provided DBTX
func (r *InstrumentRepo) CreateInstrument(ctx context.Context, dbtx db.DBTX, name string) (db.Instrument, error) {
	queries := db.New(dbtx)
	instrument, err := queries.CreateInstrument(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return db.Instrument{}, fmt.Errorf("instrument with this name already exists: %w", errors.ErrDuplicate)
		}
		return db.Instrument{}, fmt.Errorf("failed to create instrument: %w", err)
	}
	return instrument, nil
}

// GetInstrument retrieves a instrument by ID within the provided DBTX
func (r *InstrumentRepo) GetInstrument(ctx context.Context, dbtx db.DBTX, id string) (db.Instrument, error) {
	queries := db.New(dbtx)
	instrument, err := queries.GetInstrument(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Instrument{}, fmt.Errorf("instrument not found: %w", errors.ErrNotFound)
		}
		return db.Instrument{}, fmt.Errorf("failed to get instrument: %w", err)
	}
	return instrument, nil
}

// ListInstruments retrieves a paginated list of instruments within the provided DBTX
func (r *InstrumentRepo) ListInstruments(ctx context.Context, dbtx db.DBTX, limit, offset int64) ([]db.Instrument, error) {
	queries := db.New(dbtx)
	instruments, err := queries.ListInstruments(ctx, db.ListInstrumentsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list instruments: %w", err)
	}
	return instruments, nil
}

// UpdateInstrument updates an instrument within the provided DBTX
func (r *InstrumentRepo) UpdateInstrument(ctx context.Context, dbtx db.DBTX, id, name string) (db.Instrument, error) {
	queries := db.New(dbtx)
	instrument, err := queries.UpdateInstrument(ctx, db.UpdateInstrumentParams{
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

// DeleteInstrument deletes an instrument by ID within the provided DBTX
func (r *InstrumentRepo) DeleteInstrument(ctx context.Context, dbtx db.DBTX, id string) error {
	queries := db.New(dbtx)
	err := queries.DeleteInstrument(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("instrument not found: %w", errors.ErrNotFound)
		}
		return fmt.Errorf("failed to delete instrument: %w", err)
	}
	return nil
}
