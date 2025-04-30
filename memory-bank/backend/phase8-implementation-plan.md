# Phase 8 Implementation Plan: Reporting & Reconciliation

**Instructions for LLM:**
1. After completing each numbered step in the "Implementation Steps" section, wait for explicit approval from the user before proceeding to the next step.
2. After all steps in this file are completed and approved, update the memory bank to reflect the completed work for this phase.
3. Finally, use the `attempt_completion` tool to summarize the work done for this phase.

## Overview

This phase focuses on implementing initial reporting endpoints (account balance, category summary) and the logic for identifying cash reconciliation needs.

## Implementation Steps

This phase follows the **Standard Service Implementation Workflow (TDD)**, focusing on read/calculation logic:

1. **Define API (`proto/expenses/v1/expenses.proto`):** Define request/response messages for `GetAccountBalance`, `GetCategorySummary`, `ListCashDiscrepancies` RPCs. Define necessary messages for report results (e.g., `AccountBalance`, `CategorySummaryItem`, `CashDiscrepancy`). Add `service ReportingService {}` to `expenses.proto`.
2. **Define SQL (`db/queries/reporting.sql`):** Write `sqlc`-annotated queries for:
    * Calculating account balances (SUM of debits/credits for a given account).
    * Aggregating ledger entries by category (SUM of amounts for Equity splits, potentially with grouping by parent category).
    * Querying transactions and related ledger entries to identify potential cash discrepancies based on `allocation_tag` and actual cash account entries.
3. **Define Schema (`db/schema.sql`):** (Already done in Phase 0). No changes needed.
4. **Generate Code:** Run `make generate-all`. Verify generated Go code.
5. **Create/Apply Migration:** (Only if schema changed in step 3). Run `make migrate-new NAME=...` followed by `make migrate`.
6. **Write Failing Service Test (`internal/rpc/services/reporting_test.go`):**
    * Create the test file.
    * Implement `TestMain` for test database setup/teardown, seeding a variety of transactions and ledger entries across different accounts and categories to test calculations and discrepancy logic.
    * Write table-driven tests for `TestGetAccountBalance`, `TestGetCategorySummary`, and `TestListCashDiscrepancies`.
    * Test cases should cover various scenarios, including accounts with no transactions, categories with no entries, and transactions designed to create discrepancies.
    * Instantiate concrete repos (`AccountRepo`, `CategoryRepo`, `TransactionRepo`, `LedgerEntryRepo`, etc.) and the `ReportingService`.
    * Use `google/go-cmp/cmp` for assertions on calculated results and discrepancy lists.
    * Run tests and confirm they initially fail.
7. **Implement Repository Logic (`internal/repo/reporting.go`):**
    * Create the `ReportingRepo` struct.
    * Implement methods calling the `sqlc`-generated reporting queries.
8. **Implement Service Logic (`internal/rpc/services/reporting.go`):**
    * Create the `ReportingService` struct with dependencies on relevant repos.
    * Implement the `ReportingService` interface methods.
    * Call repository methods to fetch raw data.
    * Perform any necessary application-level calculations or data aggregation (e.g., handling category hierarchy for summaries, comparing allocation tags to actual cash entries for discrepancies).
    * Implement input validation (e.g., valid account/category IDs).
    * Handle errors.
    * Convert between types.
    * Continue until tests pass.
9. **Register Service (`cmd/serve.go`):** Instantiate and register all necessary repos and the `ReportingService`.
10. **Run Tests:** Run `make test`. Ensure all tests pass.
11. **Run Linter:** Run `make lint`. Fix all reported issues.
12. **Commit:** Commit changes.

## Key Considerations

* Ensure SQL queries for reporting are efficient, especially for large datasets.
* The service layer may need to perform additional logic on top of raw SQL results (e.g., hierarchical category summaries, discrepancy comparison).
* Thorough testing with diverse data scenarios is crucial to validate calculation and discrepancy logic.
* Define clear criteria for identifying cash discrepancies.

## Implementation Order

1. Define Protobuf messages and service for Reporting.
2. Define SQL queries for account balance, category summary, and discrepancy identification.
3. Run code generation (`make generate-all`).
4. Write comprehensive failing service tests in `reporting_test.go`, including tests for calculations and discrepancies.
5. Implement `ReportingRepo` methods.
6. Implement `ReportingService` methods, including application-level logic for calculations and discrepancies, until tests pass.
7. Register `ReportingService` and necessary repos in `cmd/serve.go`.
8. Run `make test` and verify all tests pass.
9. Run `make lint` and fix issues.
10. Commit changes.
