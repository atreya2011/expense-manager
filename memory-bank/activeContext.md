# Active Context: Expense Manager

## Current Work Focus

Phase 0 (Project Setup & Foundation) is **complete**. Phase 1 has begun with **Users Service implementation** now also **complete**. The focus is continuing with **Phase 1: Master Data Services**.

**Current Task:** Continue with **Phase 1: Master Data Services** by implementing the **Instruments Service**. This involves the initial steps of defining the Protobuf messages/service, creating SQL queries, generating code, and *then* proceeding with the TDD cycle (writing tests, implementing service/repo logic).

## Recent Changes (Outcomes of Phase 0)

1. **Project Structure:** Standard structure created.
2. **Configuration:** Buf, Air, sqlc, Cobra, `github.com/caarlos0/env/v11`, slog, Atlas (`schema.sql` diffing), `golangci-lint` configured.
3. **Database:** SQLite schema defined (`db/schema.sql`) including `users` and `account_users` tables. Atlas applied initial schema/seeds. DB pool established.
4. **API:** Initial proto files created. Code generation pipeline functional. Basic server starts.
5. **Tooling:** `Makefile` updated for Atlas. Air configured. Linter mandatory.
6. **Testing:** `google/go-cmp/cmp` added. Testing approach standardized.
7. **Terminology:** "Owner" renamed to "User" across schema and planned code.

## Next Steps

Follow the interleaved implementation plan outlined in `memory-bank/implementationPlan.md`, starting with:

1. **Phase F0: Frontend Setup & Tooling Integration:**
    * Verify existing frontend scaffolding (`frontend/`).
    * Install frontend dependencies (TanStack Router, TanStack Query, Tailwind CSS, Connect-ES client libraries) using Bun.
    * Configure Tailwind CSS (`tailwind.config.js`, integrate into CSS).
    * Ensure the Connect-ES client SDK is being generated correctly by `make generate-all` into `frontend/gen/client`.
    * Set up a basic API client instance using the generated SDK.
    * Create initial route structure using TanStack Router (`frontend/src/routes/`).
    * Implement a basic "Hello World" or index page to verify setup.
    * Configure frontend linting and formatting (e.g., ESLint, Prettier) if not already present.
2. **Phase F1: User Management UI (Frontend):** Implement the frontend UI for the already completed backend User Service.
3. **Phase 2: Instruments Service (Backend):** Implement the backend Instruments Service.
4. **Phase F2: Master Data Management UI (Frontend - Instruments):** Implement the frontend UI for the Instruments Service.
5. Continue alternating backend and frontend phases as per `memory-bank/implementationPlan.md`.

## Active Decisions and Considerations

1. **TDD Workflow**: Test-first development standard.
2. **Direct Database Connectivity**: Using concrete repo implementations rather than interfaces.
3. **Atlas `schema.sql` Migrations**: Schema evolution via `db/schema.sql` edits and Atlas diff/apply.
4. **DEB Implementation**: Final schema applies.
5. **Configuration**: Use `github.com/caarlos0/env/v11`.
6. **Testing Strategy**: Table-driven tests with TestMain for setup/teardown.
7. **User Concept**: Replaces "Owner" for account association.

## Important Patterns and Preferences

1. **Test Setup**: Centralized test database initialization in TestMain.
2. **Transaction Management**: Services orchestrate DB transactions (`sql.Tx`).
3. **Error Handling**: Use Connect errors; structured logging; test error paths.
4. **Code Organization**: Follow established structure. Test files alongside implementation.
5. **DRY Testing**: Common test helpers, assertions, and table-driven tests.
6. **Test Data**: Creation through sqlc-generated repository functions.

## Learnings and Project Insights

1. **Schema Precision**: Importance reinforced.
2. **Tooling Configuration**: Essential for smooth workflow.
3. **Declarative Migrations (Atlas)**: Workflow with SQL-based schema management proven effective.
4. **Mandatory Linting**: Enforces quality.
5. **TDD Setup**: Requires integration into workflow.
6. **Terminology**: Adapting terms like Owner->User early avoids confusion later.
7. **Service Implementation Pattern**: User service implementation validates the direct-to-repository approach.
8. **Test Structure**: TestMain pattern with helper functions proven effective for robust test coverage.
9. **Error Handling**: Connect error wrapping provides consistent client-friendly error responses.

## Learnings from User Service Implementation

1. The direct service-to-repository approach is effective and simplifies the data flow.
2. The TestMain pattern for test database setup and teardown works well for service-level tests.
3. Table-driven tests with `google/go-cmp/cmp` provide clear and comprehensive test coverage.
4. Registering the service handler in `cmd/serve.go` is straightforward.
