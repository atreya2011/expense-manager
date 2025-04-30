package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/atreya2011/expense-manager/internal/config"
	"github.com/atreya2011/expense-manager/internal/log"
	"github.com/atreya2011/expense-manager/internal/repo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	RunE:  runSeedCmd,
}

func init() {
	rootCmd.AddCommand(seedCmd)
}

func runSeedCmd(cmd *cobra.Command, args []string) error {
	// Initialize logger
	logger := log.NewLogger()
	if verboseMode {
		logger.Info("Verbose mode enabled")
	}
	logger.Info("Starting database seeding...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		return err
	}

	// Initialize database connection
	logger.Info("Connecting to database...", "path", cfg.Database.Path)
	db, err := repo.OpenDB(cfg.Database.Path)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return err
	}
	logger.Info("Database connection established")

	// Read seed data SQL file
	seedFilePath := "db/seed_data.sql"
	logger.Info("Reading seed data file", "path", seedFilePath)

	content, err := os.ReadFile(seedFilePath)
	if err != nil {
		logger.Error("Failed to read seed data file", "error", err)
		return err
	}

	// Execute seed statements in a transaction
	logger.Info("Executing seed data statements")
	tx, err := db.Begin()
	if err != nil {
		logger.Error("Failed to begin transaction", "error", err)
		return err
	}

	// Split the file content into separate SQL statements
	statements := strings.Split(string(content), ";")

	// Execute each statement
	for i, statement := range statements {
		stmt := strings.TrimSpace(statement)
		if stmt == "" {
			continue
		}

		if verboseMode {
			logger.Info("Executing statement", "statement_number", i+1)
		}

		_, err = tx.Exec(stmt)
		if err != nil {
			logger.Error("Failed to execute seed statement",
				"statement_number", i+1,
				"error", err)

			// Attempt to rollback transaction on error
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Error("Failed to rollback transaction", "error", rbErr)
			}

			return fmt.Errorf("failed to execute seed statement %d: %w", i+1, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit transaction", "error", err)
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		logger.Error("Failed to close database connection", "error", err)
	}

	logger.Info("Database seeding completed successfully")
	return nil
}
