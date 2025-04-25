# Expense Manager

A double-entry bookkeeping expense management system built with Go, ConnectRPC, SQLite, and sqlc.

## Prerequisites

- Go 1.21+
- Buf v2
- Air (for live reload)
- sqlc
- Atlas (for database migrations)
- SQLite3

## Development

1. Install dependencies:

```bash
go mod tidy
```

2. Generate code:

```bash
make generate
```

3. Run migrations:

```bash
make migrate-up
```

4. Run the server:

```bash
make air
```

## Project Structure

- `api/proto/`: Protocol buffer definitions
- `cmd/`: Command-line interface code
- `db/`: Database migrations and queries
- `internal/`: Internal packages
  - `auth/`: Authentication logic
  - `clock/`: Time utilities
  - `config/`: Configuration management
  - `log/`: Logging utilities
  - `repo/`: Database repositories
  - `rpc/`: RPC services
