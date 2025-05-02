# Progress: Expense Manager

## What Works (Completed in Phase 0)

Phase 0, establishing the project's foundation, is complete, reflecting the final architectural and tooling decisions. The following components are fully set up and operational:

1. **Project Structure:**
    * Standardized directory structure implemented.
    * Organization supports direct service-to-repo interaction.
2. **Database Schema & Migrations (Atlas with `schema.sql`):**
    * Canonical schema defined in `db/schema.sql`. Includes `users` table (formerly owners), `account_users` join table, and other core tables (`account_types`, `instruments`, `currencies`, `institutions`, `accounts`, `categories`, `transactions`, `ledger_entries`).
    * Atlas configured (via Makefile parameters) for schema diffing against `db/schema.sql` to generate versioned `UP` migrations (`db/migrations/`).
    * Initial schema and master data (including 'Current Year Earnings' Equity account) applied via Atlas.
3. **API Definition & Generation:**
    * Protocol Buffer definitions structured under `api/proto/expenses/v1/`.
    * Initial service definition (`service ExpenseService {}`) exists.
    * Buf v2 (`buf.yaml`, `buf.gen.yaml`) configured. `make generate` successful. Go stubs in `internal/rpc/gen/`.
4. **Data Access Generation (sqlc):**
    * sqlc (`sqlc.yaml`, `db/queries/`) configured. `make generate` successful. Go code in `internal/repo/gen/`.
5. **Core Application Setup:**
    * Go module (`go.mod`) initialized with dependencies.
    * Configuration (`internal/config/`) uses `github.com/caarlos0/env/v11`. `.env` loading via `godotenv` for local dev.
    * Structured logging (`internal/log/`) setup complete.
    * Database connection pool (`internal/repo/pool.go`) established.
    * Basic Cobra CLI (`cmd/`) structure in place.
    * Main application entry point (`main.go`) wires basic components.
6. **Development Environment:**
    * `Makefile` provides standard targets.
    * Air (`.air.toml`) configured for live reload.
    * **golangci-lint** (`.golangci.yaml`) configuration is **mandatory**.
    * Basic server starts successfully (`make run`/`make air`).
7. **Testing Foundation:**
    * Standard `testing` package will be used.
    * `google/go-cmp/cmp` added as a dependency.
8. **Database Interaction:**
    * `sqlx` integrated for enhanced database operations.
9. **Observability:**
    * OpenTelemetry tracing setup complete.

## What's Left to Build (TDD Approach)

Development proceeds phase by phase, using TDD focused on service-level tests:

1. **Phase 1: Master Data Services (Users first, then others)**
    * **Initial Setup (per service):** Define proto messages/service, define SQL queries, run `make generate-all`.
    * **TDD Cycle (per service):**
        * Write Tests: Create service tests (`*_test.go`) using standard `testing` and `cmp`.
        * Implement: Implement repo, implement service until tests pass. Register service. Check linting (`make lint`).
    * Repeat initial setup and TDD cycle for Instruments, Currencies, Institutions, AccountTypes services.
2. **Phase 2: Categories Service** (Follow setup & TDD cycle)
3. **Phase 3: Accounts Service** (Follow setup & TDD cycle - includes user linking logic via `account_users`)
4. **Phase 4: Transactions & Ledger Entries Service** (Follow TDD cycle - critical DEB validation tests)
5. **Phase 5: Reporting & Reconciliation** (Follow TDD cycle)

## Current Status

**Phase 0 (Backend Project Setup & Foundation) is complete.** Project foundation is established with the specified tools and architecture.

**Phase 1 (Backend Users Service implementation) is complete.** The backend User Service is fully implemented and tested.

**Frontend planning (Phase F0-F8) is complete.** The plan for implementing the frontend using React, TanStack Router, TanStack Query, Connect-ES, and Tailwind has been defined and integrated into the overall implementation plan.

**Next Steps:** Begin the interleaved implementation phases as outlined in `memory-bank/implementationPlan.md`, starting with **Phase F0: Frontend Setup & Tooling Integration**.

## Known Issues

* None identified specific to Phase 0 setup.

## Evolution of Project Decisions

1. **Architecture Choice**: Direct service-to-repository implementation.
2. **Database Decisions**: SQLite. Atlas with `schema.sql` diffing for migrations (UP only). DEB schema finalized (A/L/E accounts, separate Categories, Equity balancing).
3. **API Design**: ConnectRPC via Protobuf.
4. **Development Workflow**: Buf, sqlc, Air, Cobra. TDD mandatory.
5. **Configuration**: `github.com/caarlos0/env/v11` for environment variable loading.
6. **Testing**: Standard `testing` package + `google/go-cmp/cmp`.
7. **Linting**: `golangci-lint` mandatory.
8. **Terminology**: "Owner" concept renamed to "User". `owners` table becomes `users`. `account_owners` becomes `account_users`.
