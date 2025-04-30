# Phase 6 Implementation Plan: Accounts Service

**Instructions for LLM:**
1. After completing each numbered step in the "Implementation Steps" section, wait for explicit approval from the user before proceeding to the next step.
2. After all steps in this file are completed and approved, update the memory bank to reflect the completed work for this phase.
3. Finally, use the `attempt_completion` tool to summarize the work done for this phase.

## Overview

This phase focuses on implementing CRUD functionality for the `accounts` table (Asset, Liability, Equity types only) and managing the many-to-many relationship with `users` via the `account_users` join table.

## Implementation Steps

This phase follows the **Standard Service Implementation Workflow (TDD)**, with additional complexity for user linking and fetching related data:

1. **Define API (`proto/expenses/v1/expenses.proto`):** Define `Account` message (including linked `AccountType`, `Instrument`, `Institution`, `Currency`, and list of associated `User` IDs), and request/response messages for `CreateAccount`, `GetAccount`, `ListAccounts`, `UpdateAccount`, `DeleteAccount` RPCs. Add `service AccountService {}` to `expenses.proto`.
2. **Define SQL (`db/queries/account.sql`, `db/queries/account_user.sql`):** Write `sqlc`-annotated queries for Account CRUD. Also, write queries for the `account_users` table: `LinkUserToAccount`, `UnlinkUserFromAccount`, `ListUsersForAccount`. Queries for `GetAccount` and `ListAccounts` should join or fetch related data (AccountType, Instrument, etc.) and associated User IDs.
3. **Define Schema (`db/schema.sql`):** (Already done in Phase 0, includes `accounts` and `account_users` tables). No changes needed unless schema adjustments are required.
4. **Generate Code:** Run `make generate-all`. Verify generated Go code.
5. **Create/Apply Migration:** (Only if schema changed in step 3). Run `make migrate-new NAME=...` followed by `make migrate`.
6. **Write Failing Service Test (`internal/rpc/services/account_test.go`):**
    * Create the test file.
    * Implement `TestMain` for test database setup/teardown, seeding necessary master data (AccountTypes, Currencies, etc.) and some users.
    * Write table-driven tests for all CRUD operations, including linking/unlinking users during creation/update.
    * Include tests for fetching accounts and verifying that related data and associated user IDs are correctly returned.
    * Test cases for invalid inputs (e.g., linking non-existent users, invalid account types).
    * Instantiate concrete repos (`AccountRepo`, `UserRepo`, etc.) and the `AccountService`.
    * Use `google/go-cmp/cmp` for assertions.
    * Run tests and confirm they initially fail.
7. **Implement Repository Logic (`internal/repo/account.go`, `internal/repo/account_user.go`):**
    * Create `AccountRepo` and `AccountUserRepo` structs.
    * Implement Account CRUD methods in `AccountRepo`, calling `sqlc`-generated functions.
    * Implement `LinkUserToAccount`, `UnlinkUserFromAccount`, `ListUsersForAccount` in `AccountUserRepo`.
    * Queries for Get/List in `AccountRepo` should fetch related data.
8. **Implement Service Logic (`internal/rpc/services/account.go`):**
    * Create the `AccountService` struct with dependencies on relevant repos (`AccountRepo`, `AccountUserRepo`, `UserRepo`, etc.).
    * Implement the `AccountService` interface methods.
    * Call repository methods.
    * **Transaction Management:** `CreateAccount` and `UpdateAccount` methods that involve linking/unlinking users should use database transactions (`sql.Tx`) to ensure atomicity. The service layer orchestrates the transaction, passing the transaction object to repository methods.
    * Implement input validation (e.g., valid account type, existence of linked entities and users).
    * Handle errors.
    * Convert between types, populating related fields and user IDs in the `Account` message for Get/List responses.
    * Continue until tests pass.
9. **Register Service (`cmd/serve.go`):** Instantiate and register all necessary repos and the `AccountService`.
10. **Run Tests:** Run `make test`. Ensure all tests pass.
11. **Run Linter:** Run `make lint`. Fix all reported issues.
12. **Commit:** Commit changes.

## Key Considerations

* Properly manage database transactions in the service layer for operations involving multiple table modifications (e.g., creating an account and linking users).
* Ensure efficient fetching of related data (AccountType, Instrument, etc.) and associated user IDs for Get and List operations.
* Implement robust validation for linked entities and users.
* Handle potential errors when linking/unlinking users (e.g., user not found).
* Implement pagination for `ListAccounts`.

## Implementation Order

1. Define Protobuf messages and service for Accounts (including related fields and user IDs).
2. Define SQL queries for Account CRUD and `account_users` operations.
3. Run code generation (`make generate-all`).
4. Write failing service tests in `account_test.go`, including tests for user linking and fetching related data.
5. Implement `AccountRepo` and `AccountUserRepo` methods.
6. Implement `AccountService` methods, including transaction management, validation, and data fetching/conversion, until tests pass.
7. Register `AccountService` and necessary repos in `cmd/serve.go`.
8. Run `make test` and verify all tests pass.
9. Run `make lint` and fix issues.
10. Commit changes.
