# Implementation Plan: Expense Manager

**Goal:** Develop a full-stack expense management application with a Go backend API service and a React frontend, ensuring data accuracy via double-entry bookkeeping (DEB) and following a strict Test-Driven Development (TDD) methodology where applicable.

**Target Project Structure:** Adheres to the established project structure, using directories like `proto/`, `cmd/`, `db/`, `internal/`, and `frontend/`.

**Tools:**
* **Backend:** Go (latest stable), Buf v2, Air (for live reload), sqlc, **Atlas** (for schema migrations via db/schema.sql and Makefile parameters), Cobra, **golangci-lint** (mandatory), OpenTelemetry (for tracing).
* **Frontend:** Vite (build tool), Bun (package manager), ESLint, Prettier (linting/formatting).

**Libraries:**
* **Backend:** Ensure usage of the latest stable versions: connectrpc.com/connect, google.golang.org/protobuf, github.com/mattn/go-sqlite3, **github.com/jmoiron/sqlx**, OpenTelemetry libraries (e.g., `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/sdk`, exporters), **ariga.io/atlas-go-sdk** (if needed for programmatic Atlas interaction), github.com/sqlc-dev/sqlc, github.com/spf13/cobra, **github.com/caarlos0/env/v11** (for config), github.com/joho/godotenv, **log/slog** (standard library logger), **github.com/google/go-cmp/cmp** (for testing).
* **Frontend:** React, React DOM, TanStack Router, TanStack Query, @tanstack/react-query, @connectrpc/connect, @connectrpc/connect-web, @bufbuild/protobuf, Tailwind CSS.

**Methodology:**
* **Backend (TDD Focused):** Test-Driven Development (TDD) is mandatory for service-level logic using the standard `testing` package and `google/go-cmp/cmp`. Tests often interact with a real test database. Direct Service-to-Repository interaction. Declarative Migrations (Atlas) via `db/schema.sql`. Mandatory Linting (`make lint`).
* **Frontend:** Component-Based Development with React. Data fetching, caching, and state management with TanStack Query. Type safety via TypeScript and generated Connect-ES client code. Responsive Design with Tailwind CSS.
* **Interleaved Development:** Backend implementation for a service is followed by the corresponding frontend implementation before moving to the next service.

**Final Agreed ERD Summary:**
* **Core Tables:** `transactions` (Journal), `ledger_entries` (Ledger).
* **Account Structure:** `accounts` table for Asset, Liability, Equity types only. `name` is UNIQUE identifier. `account_type_id`, `instrument_id`, `institution_id`, `currency_id` provide context.
* **Ownership:** Many-to-Many via `account_users`.
* **Classification:** Separate `categories` table (Income/Expense types, optional `parent_id`). `ledger_entries` link to `categories.id` (nullable) to classify the impact, especially on Equity splits.
* **Master Data:** `account_types` (A, L, E only), `users`, `instruments`, `currencies`, `institutions`.
* **Reconciliation:** Handled by application logic using `transactions.allocation_tag` and ledger entries for cash accounts.

---

## Standard Backend Service Implementation Workflow (TDD)

(Apply this workflow to backend phases)

1. **Define/Modify API (`.proto`):** Define messages, requests, responses, and service RPCs. Register service in main `expenses.proto`.
2. **Define/Modify SQL (`db/queries/*.sql`):** Write necessary `sqlc`-annotated CRUD or specific queries.
3. **Define/Modify Schema (`db/schema.sql`):** **Only if** database structure changes are required for this service.
4. **Generate Code:** `make generate-all`. Verify generated Go code in `internal/rpc/gen/` and `internal/repo/gen/`.
5. **Create/Apply Migration:** **Only if** schema changed in step 3. Run `make migrate-new NAME=...` followed by `make migrate`.
6. **Write Failing Service Test (`internal/rpc/services/*_test.go`):**
    * Create the test file.
    * Use `TestMain` for setting up/tearing down a test database connection (potentially seeding necessary data).
    * Write table-driven tests for the new RPC.
    * Instantiate the *concrete* repository and the service implementation within the test.
    * Use `google/go-cmp/cmp` for assertions.
    * The initial test run for the new RPC should fail.
7. **Implement Repository Logic (`internal/repo/*.go`):**
    * Create the repository struct (e.g., `type UserRepo struct { db *db.Queries }`).
    * Implement methods that directly use the `sqlc`-generated `db.Queries` struct (passed during instantiation). Handle context propagation.
8. **Implement Service Logic (`internal/rpc/services/*.go`):**
    * Create the service struct with a dependency on the concrete repository.
    * Implement the service interface methods.
    * Call repository methods.
    * Perform necessary data transformations (DB model to Proto model).
    * Implement business logic and validation.
    * Handle errors using `connect.NewError`.
    * Continue until tests written in step 6 pass.
9. **Register Service (`cmd/serve.go`):** Instantiate the concrete repository, instantiate the service implementation (injecting the repo), and register the service handler with the Connect server.
10. **Run Tests:** Run `make test`. Ensure all tests in the project pass.
11. **Run Linter:** Run `make lint`. **Fix all reported issues (Mandatory).**
12. **Commit:** Commit changes.

---

## Standard Frontend Service Implementation Workflow

(Apply this workflow to frontend phases)

1. **Define Frontend Routes:** Define the necessary routes for the service's UI using TanStack Router (`frontend/src/routes/`).
2. **Integrate API Client:** Use the generated Connect-ES client SDK (`frontend/gen/client`) to interact with the backend service.
3. **Implement Data Fetching:** Use TanStack Query hooks (`useQuery`, `useQueries`, `useInfiniteQuery`) to fetch data from the backend API.
4. **Implement Mutations:** Use TanStack Query mutations (`useMutation`) for creating, updating, and deleting data via the backend API.
5. **Build UI Components:** Create React components (`frontend/src/components/`) for displaying data, forms, and other UI elements.
6. **Style Components:** Apply styling using Tailwind CSS classes.
7. **Implement Frontend Validation:** Add client-side validation to forms to provide immediate feedback to the user.
8. **Handle Loading and Error States:** Use TanStack Query's built-in states to show loading indicators and error messages.
9. **Ensure Type Safety:** Leverage TypeScript and the generated client types throughout the frontend implementation.
10. **Test (Optional but Recommended):** Write frontend tests (e.g., using React Testing Library) for components and data fetching logic.
11. **Lint and Format:** Run frontend linters (ESLint) and formatters (Prettier).
12. **Commit:** Commit changes.

---

## Combined Implementation Phases

Development proceeds phase by phase, alternating between backend and frontend implementation for each functional area.

### Phase 0: Project Setup & Foundation (Backend)

* **Goal:** Initialize the backend project workspace, configure core tools, define the database schema, apply initial migrations, and create the basic runnable server structure.
* **Status:** **Complete.**
* **Tasks:** (See detailed steps in `memory-bank/backend/phase0-implementation-plan.md`)

### Phase F0: Frontend Setup & Tooling Integration

*   **Goal:** Configure the frontend environment, install necessary dependencies, set up Tailwind CSS, and integrate the generated Connect-ES client SDK.
*   **Status:** **Partially Complete.**
*   **Tasks:**
    1. Verify existing frontend scaffolding (`frontend/`).
    2. Install frontend dependencies (TanStack Router, TanStack Query, Tailwind CSS, Connect-ES client libraries) using Bun. **(Missing TanStack Query and Connect-ES dependencies)**
    3. Configure Tailwind CSS (`tailwind.config.js`, integrate into CSS).
    4. Ensure the Connect-ES client SDK is being generated correctly by `make generate-all` into `frontend/gen/client`.
    5. Set up a basic API client instance using the generated SDK.
    6. Create initial route structure using TanStack Router (`frontend/src/routes/`).
    7. Implement a basic "Hello World" or index page to verify setup.
    8. Configure frontend linting and formatting (e.g., ESLint, Prettier) if not already present.

### Phase 1: Users Service (Backend)

* **Goal:** Implement CRUD functionality for the `users` table.
* **Status:** **Complete.**
* **Tasks:** (See detailed steps in `memory-bank/backend/phase1-implementation-plan.md`)

### Phase F1: User Management UI (Frontend)

* **Goal:** Build UI components and pages for managing users (listing, viewing, creating, updating, deleting).
* **Status:** **Planning Complete.**
* **Tasks:**
    1. Define frontend routes for user-related pages (e.g., `/users`, `/users/new`, `/users/:userId`).
    2. Use TanStack Query hooks to fetch user data via the generated `UserService` client.
    3. Create React components for displaying user lists, user details, and forms for creating/editing users.
    4. Implement mutations using TanStack Query for creating, updating, and deleting users.
    5. Style components using Tailwind CSS.
    6. Add frontend validation for user input.

### Phase 2: Instruments Service (Backend)

* **Goal:** Implement read/list functionality (CRUD if needed) for `instruments`.
* **Status:** **Pending.**
* **Tasks:** Follow the **Standard Backend Service Implementation Workflow (TDD)**. (See detailed steps in `memory-bank/backend/phase2-implementation-plan.md`)

### Phase F2: Master Data Management UI (Frontend - Instruments)

* **Goal:** Build UI for managing Instruments.
* **Status:** **Planning Complete.**
* **Tasks:**
    1. Define routes for Instruments (e.g., `/instruments`).
    2. Use TanStack Query hooks and mutations for interacting with the `InstrumentService`.
    3. Create reusable components for listing and managing Instruments.
    4. Style using Tailwind.

### Phase 3: Currencies Service (Backend)

* **Goal:** Implement read/list functionality for `currencies`.
* **Status:** **Pending.**
* **Tasks:** Follow the **Standard Backend Service Implementation Workflow (TDD)**. (See detailed steps in `memory-bank/backend/phase3-implementation-plan.md`)

### Phase F3: Master Data Management UI (Frontend - Currencies)

* **Goal:** Build UI for managing Currencies.
* **Status:** **Planning Complete.**
* **Tasks:**
    1. Define routes for Currencies (e.g., `/currencies`).
    2. Use TanStack Query hooks and mutations for interacting with the `CurrencyService`.
    3. Create reusable components for listing and managing Currencies.
    4. Style using Tailwind.

### Phase 4: Institutions Service (Backend)

* **Goal:** Implement CRUD functionality for `institutions`.
* **Status:** **Pending.**
* **Tasks:** Follow the **Standard Backend Service Implementation Workflow (TDD)**. (See detailed steps in `memory-bank/backend/phase4-implementation-plan.md`)

### Phase F4: Master Data Management UI (Frontend - Institutions)

* **Goal:** Build UI for managing Institutions.
* **Status:** **Planning Complete.**
* **Tasks:**
    1. Define routes for Institutions (e.g., `/institutions`).
    2. Use TanStack Query hooks and mutations for interacting with the `InstitutionService`.
    3. Create reusable components for listing and managing Institutions.
    4. Style using Tailwind.

### Phase 5: Categories Service (Backend)

* **Goal:** Implement CRUD for `categories`, handling optional hierarchy via `parent_id`.
* **Status:** **Pending.**
* **Tasks:** Follow the **Standard Backend Service Implementation Workflow (TDD)**. (See detailed steps in `memory-bank/backend/phase5-implementation-plan.md`)

### Phase F5: Categories Management UI (Frontend)

* **Goal:** Build UI for managing Categories, including handling the hierarchy.
* **Status:** **Planning Complete.**
* **Tasks:**
    1. Define routes for Categories (e.g., `/categories`).
    2. Use TanStack Query hooks and mutations for interacting with the `CategoryService`.
    3. Create components for listing and managing Categories, potentially visualizing the hierarchy.
    4. Style using Tailwind.

### Phase 6: Accounts Service (Backend)

* **Goal:** Implement CRUD for `accounts` (Asset, Liability, Equity only), managing users via the `account_users` join table.
* **Status:** **Pending.**
* **Tasks:** Follow the **Standard Backend Service Implementation Workflow (TDD)**. (See detailed steps in `memory-bank/backend/phase6-implementation-plan.md`)

### Phase F6: Accounts Management UI (Frontend)

* **Goal:** Build UI for managing financial accounts, including linking users.
* **Status:** **Planning Complete.**
* **Tasks:**
    1. Define routes for accounts (e.g., `/accounts`, `/accounts/new`, `/accounts/:accountId`).
    2. Use TanStack Query to fetch account data, including associated users and related master data.
    3. Create components for displaying accounts and forms for creation/editing.
    4. Implement UI for linking/unlinking users to accounts, using mutations.
    5. Style using Tailwind.

### Phase 7: Transactions & Ledger Entries Service (Backend)

* **Goal:** Implement core DEB logic for creating/reading transactions and their associated ledger entries.
* **Status:** **Pending.**
* **Tasks:** Follow the **Standard Backend Service Implementation Workflow (TDD)**. (See detailed steps in `memory-bank/backend/phase7-implementation-plan.md`)

### Phase F7: Transaction Entry & Management UI (Frontend)

* **Goal:** Build UI for creating, viewing, and listing transactions and their ledger entries.
* **Status:** **Planning Complete.**
* **Tasks:**
    1. Define routes for transactions (e.g., `/transactions`, `/transactions/new`, `/transactions/:transactionId`).
    2. Create a form for entering new transactions, including selecting accounts, categories (for Equity splits), and specifying debit/credit amounts.
    3. Implement validation in the frontend to assist users in creating balanced transactions.
    4. Use TanStack Query mutations to submit new transactions to the backend.
    5. Create components to display transaction lists and details, including associated ledger entries.
    6. Style using Tailwind.

### Phase 8: Reporting & Reconciliation (Backend)

* **Goal:** Implement initial reporting endpoints (account balance, category summary) and the logic for identifying cash reconciliation needs.
* **Status:** **Pending.**
* **Tasks:** Follow the **Standard Backend Service Implementation Workflow (TDD)**. (See detailed steps in `memory-bank/backend/phase8-implementation-plan.md`)

### Phase F8: Reporting UI (Frontend)

* **Goal:** Build UI for displaying financial reports (account balances, category summaries) and cash discrepancies.
* **Status:** **Planning Complete.**
* **Tasks:**
    1. Define routes for reporting (e.g., `/reports/account-balances`, `/reports/category-summary`, `/reports/cash-discrepancies`).
    2. Use TanStack Query hooks to fetch report data from the backend `ReportingService`.
    3. Create components to display the report data in a user-friendly format (tables, charts if needed).
    4. Style using Tailwind.

---

This revised plan accurately reflects the project's established tooling, methodology, and architecture, now including the frontend implementation and outlining an interleaved development approach.
