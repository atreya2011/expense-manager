package repo

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

// OpenDB creates a new database connection pool
func OpenDB(path string) (*sql.DB, error) {
	connStr := path
	
	// Add query parameters if not an in-memory database
	// For in-memory databases, the path will be ":memory:"
	if path != ":memory:" {
		connStr += "?_busy_timeout=5000"
	} else {
		// For in-memory database, we need to ensure it stays alive
		// by maintaining an open connection
		connStr = "file::memory:?cache=shared&_busy_timeout=5000"
	}
	
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1) // SQLite only supports one writer at a time
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := db.Ping(); err != nil {
		if err := db.Close(); err != nil {
			return nil, fmt.Errorf("error closing database: %w", err)
		}
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}

// CloseDB closes the database connection pool
func CloseDB(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
