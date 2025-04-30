# Phase 0 Implementation Plan: Project Setup & Foundation

**Instructions for LLM:**
1. After completing each numbered step in the "Implementation Steps" section, wait for explicit approval from the user before proceeding to the next step.
2. After all steps in this file are completed and approved, update the memory bank to reflect the completed work for this phase.
3. Finally, use the `attempt_completion` tool to summarize the work done for this phase.

## Overview

This phase focuses on initializing the project workspace, configuring core development tools (Buf, Air, sqlc, Atlas, Cobra, golangci-lint), defining the initial database schema using a declarative approach, applying the baseline migrations, and setting up the basic runnable server structure with configuration loading and structured logging.

## Implementation Steps

1. **Create Project Root & Git Repository:** Initialize the main project directory and set up Git for version control.
2. **Create Core Directories:** Establish necessary directories like `cmd/`, `db/`, `db/migrations/`, `internal/`, `internal/config/`, `internal/log/`, `internal/repo/`, `internal/rpc/gen/`, `internal/rpc/services/`, `proto/`, `proto/expenses/v1/`.
3. **Create Placeholder Files:** Create initial configuration files such as `.golangci.yaml`, `atlas.yaml`, `db/schema.sql`, `sqlc.yaml`, `buf.yaml`, `buf.gen.yaml`, `.air.toml`, `Makefile`, `main.go`, `cmd/root.go`, `cmd/serve.go`, `internal/config/config.go`, `internal/log/logger.go`, `internal/repo/pool.go`.
4. **Configure `.gitignore`:** Add entries for generated files and directories like `bin/`, `*.db*`, `atlas.sum`, `.env`.
5. **Configure Buf (`buf.yaml`, `buf.gen.yaml`):** Define the Protobuf module and configure generation targets for Go (`internal/rpc/gen`) and potentially frontend clients (`frontend/gen/client`).
6. **Initialize Protobuf (`proto/expenses/v1/expenses.proto`):** Create the main proto file and define an empty `service ExpenseService {}`.
7. **Initialize Go Module (`go.mod`):** Run `go mod init [module_path]`.
8. **Install Core Go Dependencies and CLI Tools:** Use `go get` to install necessary libraries (connectrpc, protobuf, sqlite3 driver, sqlc, cobra, env, godotenv, cmp) and ensure required CLI tools (buf, sqlc, atlas, air, golangci-lint) are installed globally or locally.
9. **Configure `Makefile`:** Define targets for `proto`, `sqlc`, `generate-all`, `build`, `run`, `dev`, `lint`, `migrate-new`, `migrate`, and `test`.
10. **Configure Air (`.air.toml`):** Set up live reload, ensuring `pre_cmd` includes `make generate-all`.
11. **Configure sqlc (`sqlc.yaml`):** Specify SQLite engine, schema path (`./db/schema.sql`), queries path (`./db/queries/`), output directory (`internal/repo/gen`), and options like `emit_pointers_for_null_types` and `emit_json_tags`. Set `emit_interface: false`.
12. **Configure Atlas (`atlas.yaml`):** Define the SQLite environment, linking to `db/schema.sql` and `db/migrations`.
13. **Define Initial Schema (`db/schema.sql`):** Write the complete declarative schema for all core tables (`account_types`, `users`, `instruments`, `currencies`, `institutions`, `accounts`, `account_users`, `categories`, `transactions`, `ledger_entries`) with constraints.
14. **Generate Initial Migration:** Create an empty development database (`sqlite3 db/dev.db .databases`) and run `make migrate-new NAME=baseline` to generate the first migration script.
15. **Apply Initial Migration:** Run `make migrate` to apply the baseline schema to the main database file (`db/expenses.db`).
16. **Seed Master Data:** Add `INSERT` statements for initial master data (e.g., `account_types`, `currencies`, 'Current Year Earnings' Equity account) either directly in `db/schema.sql` within an `atlas:seed` block or by creating a dedicated seed migration file. Run `make migrate` again if using a new migration.
17. **Configure Logging (`internal/log/logger.go`):** Implement a basic structured logger using `log/slog`.
18. **Configure Clock (`internal/clock/clock.go`):** Implement a simple clock interface and a real clock implementation.
19. **Configure Configuration (`internal/config/config.go`):** Define a configuration struct and use `github.com/caarlos0/env/v11` and `godotenv` to load settings from environment variables and a `.env` file. Create a `.env` file (and add to `.gitignore`).
20. **Configure DB Pool (`internal/repo/pool.go`):** Set up the database connection pool using the configured database path.
21. **Configure Cobra (`cmd/root.go`, `cmd/serve.go`):** Set up the basic CLI structure and the `serve` command to load configuration, initialize logger, clock, and DB pool correctly.
22. **Configure Main (`main.go`):** Wire up the root command and execute the CLI application.
23. **Configure Linting (`.golangci.yaml`):** Set up basic linting rules and enable essential linters.
24. **Initial Generation & Tidy:** Run `make generate-all`, `go mod tidy`.
25. **Test Run & Lint:** Run the server (`make air` or `make run`) to ensure it starts without errors. Run `make lint` and fix any initial linting issues.
26. **Commit:** Commit all changes with a descriptive message.

## Key Considerations

* Ensure all required CLI tools are installed and accessible in the development environment.
* The `db/schema.sql` is the single source of truth for the database schema.
* Atlas migrations are UP only.
* The `.env` file should be ignored by Git.
* Mandatory linting is enforced from this phase onwards.

## Implementation Order

1. Create project structure and initial files.
2. Configure Git and `.gitignore`.
3. Configure Buf, sqlc, Atlas, Air, Cobra, golangci-lint.
4. Initialize Go module and install dependencies/tools.
5. Define the complete database schema in `db/schema.sql`.
6. Generate and apply the baseline migration.
7. Add and apply seed data (if using a migration).
8. Implement configuration loading, logging, clock, and DB pool.
9. Set up the basic Cobra CLI and main entry point.
10. Run code generation and tidy Go modules.
11. Test server startup and run initial linting.
12. Commit changes.
