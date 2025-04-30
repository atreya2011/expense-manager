# Phase 7 Implementation Plan: Transactions & Ledger Entries Service

**Instructions for LLM:**
1. After completing each numbered step in the "Implementation Steps" section, wait for explicit approval from the user before proceeding to the next step.
2. After all steps in this file are completed and approved, update the memory bank to reflect the completed work for this phase.
3. Finally, use the `attempt_completion` tool to summarize the work done for this phase.

## Overview

This phase focuses on implementing the core Double-Entry Bookkeeping (DEB) logic for creating and reading transactions and their associated ledger entries. This is a critical phase requiring extensive testing to ensure data integrity.

## Implementation Steps

This phase follows the **Standard Service Implementation Workflow (TDD)**, with a strong emphasis on the `CreateTransaction` logic and testing:

1. **Define API (`proto/expenses/v1/expenses.proto`):** Define `Transaction` message (including `ledger_entries`), `LedgerEntry` message (including links to `Account` and optional `Category`), and request/response messages for `CreateTransaction`, `GetTransaction`, `ListTransactions`, `UpdateTransactionMetadata`, `DeleteTransaction` RPCs. Add `service TransactionService {}` to `expenses.proto`.
2. **Define SQL (`db/queries/transaction.sql`, `db/queries/ledger_entry.sql`):** Write `sqlc`-annotated queries for Transaction and Ledger Entry CRUD. Queries for `GetTransaction` and `ListTransactions` should fetch associated ledger entries and potentially related account/category data.
3. **Define Schema (`db/schema.sql`):** (Already done in Phase 0, includes `transactions` and `ledger_entries` tables). No changes needed unless schema adjustments are required.
4. **Generate Code:** Run `make generate-all`. Verify generated Go code.
5. **Create/Apply Migration:** (Only if schema changed in step 3). Run `make migrate-new NAME=...` followed by `make migrate`.
6. **Write Failing Service Test (`internal/rpc/services/transaction_test.go`):**
    * Create the test file.
    * Implement `TestMain` for test database setup/teardown, seeding necessary accounts (including the Equity account) and categories.
    * **CRITICAL TDD FOCUS:** Write comprehensive table-driven tests for `TestCreateTransaction`. Test cases must cover:
        * Successful creation with balanced debits/credits.
        * Creation with unbalanced debits/credits (should fail).
        * Creation with invalid account IDs.
        * Creation with invalid category IDs (for Equity splits).
        * Creation with correct `category_id` linking for Equity splits.
        * Ensuring database transactionality (rollback on error).
    * Write tests for `TestGetTransaction`, `TestListTransactions`, `TestUpdateTransactionMetadata`, `TestDeleteTransaction`.
    * Instantiate concrete repos (`TransactionRepo`, `LedgerEntryRepo`, `AccountRepo`, `CategoryRepo`, etc.) and the `TransactionService`.
    * Use `google/go-cmp/cmp` for assertions.
    * Run tests and confirm they initially fail.
7. **Implement Repository Logic (`internal/repo/transaction.go`, `internal/repo/ledger_entry.go`):**
    * Create `TransactionRepo` and `LedgerEntryRepo` structs.
    * Implement CRUD methods in both repos, calling `sqlc`-generated functions.
    * Methods for Get/List should fetch associated data.
8. **Implement Service Logic (`internal/rpc/services/transaction.go`):**
    * Create the `TransactionService` struct with dependencies on relevant repos.
    * Implement the `TransactionService` interface methods.
    * **`CreateTransaction` Logic:**
        * Validate that the sum of debits equals the sum of credits in the input ledger entries. Return an error if unbalanced.
        * Validate that all referenced account and category IDs exist and are valid for the entry type (e.g., category only for Equity splits).
        * Start a database transaction.
        * Call `TransactionRepo.CreateTransaction` within the transaction.
        * Iterate through the input ledger entries, calling `LedgerEntryRepo.CreateLedgerEntry` for each within the same transaction.
        * Commit the transaction if all operations succeed; rollback on any error.
    * Implement `GetTransaction`, `ListTransactions`, `UpdateTransactionMetadata`, `DeleteTransaction`, calling repository methods and handling errors/conversions. `DeleteTransaction` should also delete associated ledger entries within a transaction.
    * Handle errors using `connect.NewError`.
    * Convert between types.
    * Continue until tests pass.
9. **Register Service (`cmd/serve.go`):** Instantiate and register all necessary repos and the `TransactionService`.
10. **Run Tests:** Run `make test`. Ensure all tests pass.
11. **Run Linter:** Run `make lint`. Fix all reported issues.
12. **Commit:** Commit changes.

## Key Considerations

* The `CreateTransaction` method is the core of the DEB implementation and requires rigorous testing.
* Ensure strict validation of debit/credit balance and linked entity validity.
* Database transactionality is essential for `CreateTransaction` and `DeleteTransaction` to maintain data integrity.
* Properly handle the nullable `category_id` in `ledger_entries`.
* Implement pagination for `ListTransactions`.

## Implementation Order

1. Define Protobuf messages and service for Transactions and Ledger Entries.
2. Define SQL queries for Transaction and Ledger Entry CRUD.
3. Run code generation (`make generate-all`).
4. Write comprehensive failing service tests in `transaction_test.go`, with a strong focus on `TestCreateTransaction` and its validation/transactionality.
5. Implement `TransactionRepo` and `LedgerEntryRepo` methods.
6. Implement `TransactionService` methods, including the critical `CreateTransaction` logic with validation and transaction management, until tests pass.
7. Register `TransactionService` and necessary repos in `cmd/serve.go`.
8. Run `make test` and verify all tests pass.
9. Run `make lint` and fix issues.
10. Commit changes.
