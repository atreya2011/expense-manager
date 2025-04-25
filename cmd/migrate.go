package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/atreya2011/expense-manager/internal/config"
	"github.com/atreya2011/expense-manager/internal/log"
	"github.com/spf13/cobra"
)

var (
	// Migration flags
	migrationName string
	migrateDown   bool
	showStatus    bool
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run database migrations using Atlas. Can apply, rollback, or create new migrations.`,
	RunE:  runMigrateCmd,
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Add flags specific to the migrate command
	migrateCmd.Flags().StringVarP(&migrationName, "name", "n", "", "Name for the new migration file")
	migrateCmd.Flags().BoolVarP(&migrateDown, "down", "d", false, "Roll back the last migration")
	migrateCmd.Flags().BoolVarP(&showStatus, "status", "s", false, "Show migration status")
}

func runMigrateCmd(cmd *cobra.Command, args []string) error {
	logger := log.NewLogger()
	logger.Info("Running database migrations...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		return err
	}

	// Create the migrations directory if it doesn't exist
	migrationsDir := "./db/migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(migrationsDir, 0755); err != nil {
			logger.Error("Failed to create migrations directory", "error", err)
			return err
		}
		logger.Info("Created migrations directory", "path", migrationsDir)
	}

	// Ensure parent directory for database exists
	dbDir := filepath.Dir(cfg.Database.Path)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			logger.Error("Failed to create database directory", "error", err)
			return err
		}
		logger.Info("Created database directory", "path", dbDir)
	}

	var atlasCmd *exec.Cmd

	switch {
	case migrationName != "": // Create a new migration
		atlasCmd = exec.Command("atlas", "migrate", "new", migrationName, "--dir", migrationsDir)
		logger.Info("Creating new migration", "name", migrationName)

	case migrateDown: // Roll back the last migration
		atlasCmd = exec.Command("atlas", "migrate", "down", "--url", fmt.Sprintf("sqlite://%s", cfg.Database.Path), "--dir", fmt.Sprintf("file://%s", migrationsDir), "1")
		logger.Info("Rolling back the last migration")

	case showStatus: // Show migration status
		atlasCmd = exec.Command("atlas", "migrate", "status", "--url", fmt.Sprintf("sqlite://%s", cfg.Database.Path), "--dir", fmt.Sprintf("file://%s", migrationsDir))
		logger.Info("Checking migration status")

	default: // Apply all migrations
		atlasCmd = exec.Command("atlas", "migrate", "apply", "--url", fmt.Sprintf("sqlite://%s", cfg.Database.Path), "--dir", fmt.Sprintf("file://%s", migrationsDir))
		logger.Info("Applying all migrations")
	}

	// Connect the command's stdout and stderr to our process
	atlasCmd.Stdout = os.Stdout
	atlasCmd.Stderr = os.Stderr

	// Run the atlas command
	if err := atlasCmd.Run(); err != nil {
		logger.Error("Migration command failed", "error", err)
		return err
	}

	logger.Info("Migration command completed successfully")
	return nil
}
