# Technical Context: Expense Manager

## Technologies Used

### Core Technologies

1. **Go (Latest Stable Recommended)**: Backend language, Go Modules.
2. **ConnectRPC**: RPC framework.
3. **Protocol Buffers (proto3)**: API IDL.
4. **SQLite (v3)**: Embedded relational database.

### Development & Build Tools

1. **Buf (v2 Recommended)**: Protobuf management (`buf.yaml`, `buf.gen.yaml`). Output: `internal/rpc/gen/` (Go), `frontend/gen/client` (ES).
2. **sqlc**: Type-safe Go code generation from SQL (`sqlc.yaml`, `db/queries/`). Output: `internal/repo/gen/`.
3. **Atlas**: Declarative database schema migration tool. Manages schema via diffing against `db/schema.sql` using parameters in the Makefile.
4. **Cobra**: CLI framework (`cmd/`).
5. **Makefile**: Centralized build/dev task runner. Includes targets like `generate-all`, `migrate`, `migrate-new`, `test`, `lint`, `build`, `run`, `seed`.
6. **golangci-lint**: **Mandatory** Go static analysis runner (`.golangci.yaml`, `make lint`).
7. **richgo**: Used via `make test` for enhanced test output.

### Supporting Libraries (Ensure Latest Stable)

1. **`github.com/caarlos0/env/v11`**: Environment variable configuration loading.
2. **`github.com/joho/godotenv`**: `.env` file loading.
3. **`log/slog` (Go Standard Library)**: Structured logging implementation (`internal/log/logger.go`).
4. **`github.com/mattn/go-sqlite3`**: SQLite driver. Note: Using `jmoiron/sqlx` for enhanced database operations over standard `database/sql`.
5. **`github.com/jmoiron/sqlx`**: Database library providing extensions to `database/sql`.
6. **OpenTelemetry**: For distributed tracing. Includes relevant Go libraries (e.g., `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/sdk`, exporters).
7. **Go Standard `testing` package**: Core testing framework.
8. **`github.com/google/go-cmp/cmp`**: Struct comparisons in tests.
9. **(Atlas Go Provider if needed)**

## Development Setup

### Prerequisites

```shell
- Go (latest stable)
- Buf CLI (v2 recommended)
- Air CLI
- sqlc CLI
- Atlas CLI
- SQLite3 development libraries
- golangci-lint CLI
```

### Development Workflow (TDD Focused)

1. **Write Tests:** Define service behavior in `*_test.go` (using `testing`/`cmp`).
2. **Define/Modify API:** Edit `.proto`.
3. **Define/Modify DB Queries:** Edit `.sql`.
4. **Define/Modify DB Schema:** Edit **`db/schema.sql`**.
5. **Generate Code & Prep Migration:** `make generate-all` (buf & sqlc). `make migrate-new name=...` (atlas diff).
6. **Apply Migration:** `make migrate` (atlas apply).
7. **Implement Logic:** Write `internal/repo/` and `internal/rpc/services/` until tests pass.
8. **Run Server:** `make run` (or use external tools like `air` manually if desired).
9. **Run Tests:** `make test` (uses `richgo`).
10. **Refactor.**
11. **Lint Code:** `make lint` (Mandatory Pass).
12. **Build Binary:** `make build`.
13. **Seed Data (Optional):** `make seed`.

## Technical Constraints

1. **SQLite:** Single-writer concurrency.
2. **ConnectRPC:** Client compatibility.
3. **Go:** Manual DI, error handling.
4. **Atlas Schema Diff Migrations:** Schema evolution via `db/schema.sql`. No down migrations.
5. **Configuration:** Relies on environment variables (`github.com/caarlos0/env/v11`).

## Dependency Management

1. **Go Modules:** `go.mod` defines dependencies. Manage via `go get`. Run `make tidy`.

## Tool Usage Patterns

### Code Generation Pattern

```shell
.proto -> Buf Generate -> internal/rpc/gen/ (Go), frontend/gen/client (ES)
.sql   -> sqlc Generate -> internal/repo/gen/ (Go)
```

### Database Migration Pattern (Atlas Schema Diff)

```shell
db/schema.sql (Edit) -> Atlas Diff (make migrate-new name=...) -> db/migrations/*.sql (Generated) -> Atlas Apply (make migrate) -> db/expenses.db (Updated)
```

### Dependency Injection Pattern (Manual Constructor)

```shell
DB Pool -> Repo Impl -> Service Impl -> RPC Handler Registration
```

### Test-Driven Development Pattern (Backend)

```shell
Write Failing Service Test (*_test.go using testing/cmp) -> Implement Code (Repo + Service) -> Run Tests & Lint (make test && make lint) -> Refactor
```
