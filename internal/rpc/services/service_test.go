package services

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/atreya2011/expense-manager/internal/clock"
	"github.com/atreya2011/expense-manager/internal/repo"
	db "github.com/atreya2011/expense-manager/internal/repo/gen"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	testDBPath = ":memory:"
)

var (
	// Global test database connection
	testDB *sqlx.DB

	// Global repositories for tests
	userRepo       *repo.UserRepo
	instrumentRepo *repo.InstrumentRepo

	// Test clock for predictable timestamps
	testClock clock.Clock

	// Test logger for predictable logging
	testLogger *slog.Logger
)

// TestMain handles setup and teardown for all tests
func TestMain(m *testing.M) {
	// Parse test flags
	flag.Parse()

	// Setup test environment
	if err := setupTestEnvironment(); err != nil {
		os.Exit(1)
	}

	// Run tests
	exitCode := m.Run()

	// Cleanup test environment
	teardownTestEnvironment()

	os.Exit(exitCode)
}

// setupTestEnvironment initializes the test database and services
func setupTestEnvironment() error {
	// Create a new in-memory test database
	var err error
	testDB, err = repo.OpenDB(testDBPath)
	if err != nil {
		return err
	}

	// Create test schema
	if err := createTestSchema(testDB); err != nil {
		return err
	}

	// Initialize repositories
	userRepo = repo.NewUserRepo(testDB)
	instrumentRepo = repo.NewInstrumentRepo(testDB)

	// Initialize test clock
	testClock = clock.NewMockClock(time.Date(2025, 4, 26, 12, 0, 0, 0, time.UTC))

	// Initialize test logger that discards output
	testLogger = slog.New(slog.DiscardHandler)

	return nil
}

// teardownTestEnvironment cleans up after tests
func teardownTestEnvironment() {
	if testDB != nil {
		if err := testDB.Close(); err != nil {
			panic("Failed to close test database: " + err.Error())
		}
	}
	// No need to remove the database file as we're using in-memory database
}

// Create test schema
func createTestSchema(db *sqlx.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (email)
		)
	`)
	if err != nil {
		return err
	}

	// Create instruments table
	_, err = db.Exec(`
		CREATE TABLE instruments (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE (name)
		)
	`)
	return err
}

// resetTestDB clears all data from the test database for a fresh test
func resetTestDB(t *testing.T) {
	t.Helper()

	// Delete all data from tables
	tables := []string{"users", "instruments"}
	for _, table := range tables {
		_, err := testDB.Exec("DELETE FROM " + table)
		if err != nil {
			t.Fatalf("Failed to clear table %s: %v", table, err)
		}
	}
}

// createTestUser inserts a test user into the database using the provided DBTX
func createTestUser(t *testing.T, dbtx db.DBTX, name, email string) db.User {
	t.Helper()

	var user db.User

	ctx := context.Background()
	params := db.CreateUserParams{
		Name:  name,
		Email: email,
	}

	// Use the repository to create the user
	var err error
	user, err = userRepo.CreateUser(ctx, dbtx, params)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

// createTestInstrument inserts a test instrument into the database using the provided DBTX
func createTestInstrument(t *testing.T, dbtx db.DBTX, name string) db.Instrument {
	t.Helper()

	var instrument db.Instrument

	ctx := context.Background()

	// Use the repository to create the instrument
	var err error
	instrument, err = instrumentRepo.CreateInstrument(ctx, dbtx, name)
	if err != nil {
		t.Fatalf("Failed to create test instrument: %v", err)
	}

	return instrument
}

// assertError checks if an error matches the expected condition
func assertError(t *testing.T, err error, expectError bool, message string) {
	t.Helper()

	if expectError && err == nil {
		t.Errorf("Expected error: %s", message)
	} else if !expectError && err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
