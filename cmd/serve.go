package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atreya2011/expense-manager/internal/clock"
	"github.com/atreya2011/expense-manager/internal/config"
	"github.com/atreya2011/expense-manager/internal/log"
	"github.com/atreya2011/expense-manager/internal/repo"
	"github.com/atreya2011/expense-manager/internal/rpc/gen/expenses/v1/expensesv1connect"
	"github.com/atreya2011/expense-manager/internal/rpc/services"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the expense manager server",
	RunE:  runServeCmd,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func runServeCmd(cmd *cobra.Command, args []string) error {
	// Initialize logger
	logger := log.NewLogger()
	if verboseMode {
		// Set more verbose logging when verbose flag is enabled
		logger.Info("Verbose mode enabled")
	}
	logger.Info("Starting server...")

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

	// Check if the database schema is correct by looking for one of our tables
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='account_types'").Scan(&count)
	if err != nil {
		logger.Error("Failed to check database schema", "error", err)
		return err
	}
	if count == 0 {
		logger.Error("Database schema is incorrect. Make sure migrations have been applied.", "error", "account_types table not found")
		logger.Info("Run 'make migrate-up' to apply Atlas migrations")
		return fmt.Errorf("account_types table not found")
	}
	logger.Info("Database schema verified")

	// Initialize repositories
	userRepo := repo.NewUserRepo(db)
	logger.Info("User repository initialized")

	// Initialize clock
	clk := clock.NewRealClock()
	logger.Info("Clock initialized")

	// Initialize Connect RPC services
	userService := services.NewUserService(userRepo, clk)
	logger.Info("User service initialized")

	// Create router
	mux := http.NewServeMux()

	// Register Connect RPC services
	userPath, userHandler := expensesv1connect.NewUserServiceHandler(userService)
	mux.Handle(userPath, userHandler)
	logger.Info("User service registered", "path", userPath)

	// Configure server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Info("Server listening", "address", addr)

	// Create server with h2c for HTTP/2 without TLS
	server := &http.Server{
		Addr:         addr,
		Handler:      h2c.NewHandler(mux, &http2.Server{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
	logger.Info("Shutdown signal received, initiating graceful shutdown...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Close database connection
	if err := db.Close(); err != nil {
		logger.Error("Failed to close database connection", "error", err)
	}

	// Shutdown server
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server graceful shutdown failed", "error", err)
		return err
	}

	logger.Info("Server shutdown gracefully")
	return nil
}
