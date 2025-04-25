# Database Migrations with Atlas

This project uses [Atlas](https://atlasgo.io/) for managing database migrations.

## Migration Commands

The following commands are available for working with migrations:

```bash
# Apply all pending migrations
make migrate-up

# Revert the last migration
make migrate-down

# Check migration status
make migrate-status

# Create a new migration file
make migrate-new name=your_migration_name
```

You can also use the CLI directly:

```bash
# Apply migrations
./bin/expense-manager migrate

# Show migration status
./bin/expense-manager migrate -s

# Rollback the last migration
./bin/expense-manager migrate -d

# Create a new migration
./bin/expense-manager migrate -n your_migration_name
```

## Migration Files

Migration files are stored in the `db/migrations` directory and follow the format `YYYYMMDDhhmmss_name.sql`.

## Initial Schema

The initial migrations include:

1. `baseline` - Creates all tables and relationships (account_types, users, instruments, currencies, institutions, accounts, account_users, categories, transactions, ledger_entries)
2. `seed_data` - Populates essential reference data like account types, default currency, and basic instruments

## Migration Strategy

Atlas handles the complete migration lifecycle:

1. Tracks applied migrations in a schema_migrations table
2. Provides status checking to identify pending migrations
3. Enables rollback of migrations when needed
4. Supports both versioned migrations and declarative schema management
