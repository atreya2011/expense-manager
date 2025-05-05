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

`make dev` will reset the database, run migrations, and start the server with live reload.
(or) you can run the following commands manually:

1. Install dependencies:

```bash
go mod tidy
```

2. Generate code:

```bash
make generate-all
```

3. Run migrations:

```bash
make migrate
```

4. Seed Master Data:

```bash
make seed
```

5. Run the server:

```bash
make run
```

## Project Structure

- `cmd/`: Command-line interface code
- `db/`: Database migrations and queries
- `internal/`: Internal packages
  - `auth/`: Authentication logic
  - `clock/`: Time utilities
  - `config/`: Configuration management
  - `log/`: Logging utilities
  - `repo/`: Database repositories
  - `rpc/`: RPC services
- `proto/`: Protocol buffer definitions

## Database Migrations with Atlas

This project uses [Atlas](https://atlasgo.io/) for managing database migrations.

### Migration Commands

The following commands are available for working with migrations:

```bash
# Apply all pending migrations
make migrate

# Check migration status
make migrate-status

# Create a new migration file based on the current schema
make migrate-new name=your_migration_name
```

### Migration Files

Migration files are stored in the `db/migrations` directory and follow the format `YYYYMMDDhhmmss_name.sql`.

### Migration Strategy

Atlas handles the complete migration lifecycle:

1. Tracks applied migrations in a schema_migrations table
2. Provides status checking to identify pending migrations
3. Supports versioned migrations
