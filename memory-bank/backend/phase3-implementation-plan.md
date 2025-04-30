# Phase 3 Implementation Plan: Currencies Service

**Instructions for LLM:**
1. After completing each numbered step in the "Implementation Steps" section, wait for explicit approval from the user before proceeding to the next step.
2. After all steps in this file are completed and approved, update the memory bank to reflect the completed work for this phase.
3. Finally, use the `attempt_completion` tool to summarize the work done for this phase.

## Overview

This phase focuses on implementing read/list functionality for the `currencies` table. Full CRUD is likely not required based on the project brief, which mentions focusing on JPY initially and implies currencies are master data.

## Implementation Steps

This phase follows the **Standard Service Implementation Workflow (TDD)**, focusing on Get and List operations:

1. **Define API (`proto/expenses/v1/expenses.proto`):** Define `Currency` message, and request/response messages for `GetCurrency`, `ListCurrencies` RPCs. Add `service CurrencyService {}` to `expenses.proto`.
2. **Define SQL (`db/queries/currency.sql`):** Write `sqlc`-annotated queries for `GetCurrency` (by ID or code) and `ListCurrencies`.
3. **Define Schema (`db/schema.sql`):** (Already done in Phase 0, includes `currencies` table). No changes needed unless schema adjustments are required.
4. **Generate Code:** Run `make generate-all`. Verify generated Go code for Protobuf stubs (`internal/rpc/gen`) and sqlc repository methods (`internal/repo/gen`).
5. **Create/Apply Migration:** (Only if schema changed in step 3). Run `make migrate-new NAME=...` followed by `make migrate`.
6. **Write Failing Service Test (`internal/rpc/services/currency_test.go`):**
    * Create the test file.
    * Implement `TestMain` for setting up and tearing down a test database connection, potentially seeding initial currency data (like JPY).
    * Write table-driven tests for `TestGetCurrency` and `TestListCurrencies`.
    * Instantiate the concrete `CurrencyRepo` and `CurrencyService` implementations within the tests.
    * Use `google/go-cmp/cmp` for assertions.
    * Run tests and confirm they initially fail.
7. **Implement Repository Logic (`internal/repo/currency.go`):**
    * Create the `CurrencyRepo` struct.
    * Implement `GetCurrency` and `ListCurrencies` methods, calling the corresponding `sqlc`-generated functions.
8. **Implement Service Logic (`internal/rpc/services/currency.go`):**
    * Create the `CurrencyService` struct.
    * Implement the `CurrencyService` interface methods (`GetCurrency`, `ListCurrencies`).
    * Call the respective `CurrencyRepo` methods.
    * Implement input validation (e.g., non-empty currency code/ID).
    * Handle errors, mapping repository errors to appropriate `connect.NewError` instances.
    * Convert between database model types and Protobuf message types.
    * Continue implementation until tests pass.
9. **Register Service (`cmd/serve.go`):** Instantiate the `CurrencyRepo`, instantiate the `CurrencyService`, and register the `CurrencyService` handler.
10. **Run Tests:** Run `make test`. Ensure all tests pass.
11. **Run Linter:** Run `make lint`. Fix all reported issues.
12. **Commit:** Commit changes.

## Key Considerations

* Ensure JPY currency data is seeded in the test database for `TestMain`.
* Consider supporting retrieval by both ID and currency code in `GetCurrency`.
* Implement pagination for `ListCurrencies`.
* Focus on read operations as CRUD is likely not needed for this master data.

## Implementation Order

1. Define Protobuf messages and service for Currencies.
2. Define SQL queries for Get and List operations.
3. Run code generation (`make generate-all`).
4. Write failing service tests in `currency_test.go`.
5. Implement `CurrencyRepo` methods.
6. Implement `CurrencyService` methods until tests pass.
7. Register `CurrencyService` in `cmd/serve.go`.
8. Run `make test` and verify all tests pass.
9. Run `make lint` and fix issues.
10. Commit changes.
