# Phase 5 Implementation Plan: Categories Service

**Instructions for LLM:**
1. After completing each numbered step in the "Implementation Steps" section, wait for explicit approval from the user before proceeding to the next step.
2. After all steps in this file are completed and approved, update the memory bank to reflect the completed work for this phase.
3. Finally, use the `attempt_completion` tool to summarize the work done for this phase.

## Overview

This phase focuses on implementing CRUD functionality for the `categories` table, including handling the optional hierarchy via the `parent_id` field.

## Implementation Steps

This phase follows the **Standard Service Implementation Workflow (TDD)**:

1. **Define API (`proto/expenses/v1/expenses.proto`):** Define `Category` message (including optional `parent_id`), and request/response messages for `CreateCategory`, `GetCategory`, `ListCategories`, `UpdateCategory`, `DeleteCategory` RPCs. Add `service CategoryService {}` to `expenses.proto`.
2. **Define SQL (`db/queries/category.sql`):** Write `sqlc`-annotated queries for `CreateCategory`, `GetCategory`, `ListCategories`, `UpdateCategory`, `DeleteCategory`. Consider queries to check for child categories before deletion.
3. **Define Schema (`db/schema.sql`):** (Already done in Phase 0, includes `categories` table with `parent_id`). No changes needed unless schema adjustments are required.
4. **Generate Code:** Run `make generate-all`. Verify generated Go code.
5. **Create/Apply Migration:** (Only if schema changed in step 3). Run `make migrate-new NAME=...` followed by `make migrate`.
6. **Write Failing Service Test (`internal/rpc/services/category_test.go`):**
    * Create the test file.
    * Implement `TestMain` for test database setup/teardown.
    * Write table-driven tests for all CRUD operations, including cases with and without `parent_id`.
    * Include tests for attempting to delete a category with children or linked ledger entries.
    * Instantiate concrete repo and service.
    * Use `google/go-cmp/cmp` for assertions.
    * Run tests and confirm they initially fail.
7. **Implement Repository Logic (`internal/repo/category.go`):**
    * Create the `CategoryRepo` struct.
    * Implement CRUD methods calling `sqlc`-generated functions.
    * Add a method to check for child categories or linked ledger entries.
8. **Implement Service Logic (`internal/rpc/services/category.go`):**
    * Create the `CategoryService` struct.
    * Implement the `CategoryService` interface methods.
    * Call `CategoryRepo` methods.
    * Implement input validation (e.g., non-empty name, valid `parent_id`).
    * Add logic in `DeleteCategory` to prevent deletion if children or linked ledger entries exist, returning an appropriate error.
    * Handle errors.
    * Convert between types.
    * Continue until tests pass.
9. **Register Service (`cmd/serve.go`):** Instantiate and register the `CategoryRepo` and `CategoryService`.
10. **Run Tests:** Run `make test`. Ensure all tests pass.
11. **Run Linter:** Run `make lint`. Fix all reported issues.
12. **Commit:** Commit changes.

## Key Considerations

* Properly handle the nullable `parent_id` field in Protobuf messages, SQL queries, and Go code.
* Implement validation to prevent circular dependencies in the category hierarchy.
* Add logic to prevent deletion of categories that are in use (have children or linked ledger entries).
* Implement pagination for `ListCategories`.

## Implementation Order

1. Define Protobuf messages and service for Categories (including `parent_id`).
2. Define SQL queries for Category CRUD operations and checks for dependencies.
3. Run code generation (`make generate-all`).
4. Write failing service tests in `category_test.go`, including hierarchy and deletion constraint tests.
5. Implement `CategoryRepo` methods.
6. Implement `CategoryService` methods, including validation and deletion constraints, until tests pass.
7. Register `CategoryService` in `cmd/serve.go`.
8. Run `make test` and verify all tests pass.
9. Run `make lint` and fix issues.
10. Commit changes.
