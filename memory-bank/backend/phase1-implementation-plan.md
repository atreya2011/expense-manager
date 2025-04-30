# Phase 1 Implementation Plan: Users Service

**Instructions for LLM:**
1. After completing each numbered step in the "Implementation Steps" section, wait for explicit approval from the user before proceeding to the next step.
2. After all steps in this file are completed and approved, update the memory bank to reflect the completed work for this phase.
3. Finally, use the `attempt_completion` tool to summarize the work done for this phase.

## Overview

This phase focuses on implementing the full CRUD (Create, Read, Update, Delete) functionality for the `users` table, following the established Test-Driven Development (TDD) workflow and architectural patterns.

## Implementation Steps

This phase follows the **Standard Service Implementation Workflow (TDD)**:

1. **Define API (`proto/expenses/v1/expenses.proto`):** Define `User` message, and request/response messages for `CreateUser`, `GetUser`, `ListUsers`, `UpdateUser`, `DeleteUser` RPCs. Add `service UserService {}` to `expenses.proto`.
2. **Define SQL (`db/queries/user.sql`):** Write `sqlc`-annotated queries for `CreateUser`, `GetUser`, `ListUsers`, `UpdateUser`, `DeleteUser`.
3. **Define Schema (`db/schema.sql`):** (Already done in Phase 0, includes `users` table). No changes needed unless schema adjustments are required.
4. **Generate Code:** Run `make generate-all`. Verify generated Go code for Protobuf stubs (`internal/rpc/gen`) and sqlc repository methods (`internal/repo/gen`).
5. **Create/Apply Migration:** (Only if schema changed in step 3). Run `make migrate-new NAME=...` followed by `make migrate`.
6. **Write Failing Service Test (`internal/rpc/services/user_test.go`):**
    * Create the test file.
    * Implement `TestMain` for setting up and tearing down a test database connection.
    * Write table-driven tests for each RPC (`TestCreateUser`, `TestGetUser`, `TestListUsers`, `TestUpdateUser`, `TestDeleteUser`).
    * Instantiate the concrete `UserRepo` and `UserService` implementations within the tests.
    * Use `google/go-cmp/cmp` for comparing expected and actual results.
    * Run tests and confirm they initially fail due to missing implementation.
7. **Implement Repository Logic (`internal/repo/user.go`):**
    * Create the `UserRepo` struct with a dependency on the `sqlc`-generated `db.Queries`.
    * Implement `CreateUser`, `GetUser`, `ListUsers`, `UpdateUser`, `DeleteUser` methods, calling the corresponding `sqlc`-generated functions. Handle `context.Context` propagation.
8. **Implement Service Logic (`internal/rpc/services/user.go`):**
    * Create the `UserService` struct with a dependency on the concrete `UserRepo`.
    * Implement the `UserService` interface methods (`CreateUser`, `GetUser`, `ListUsers`, `UpdateUser`, `DeleteUser`).
    * Call the respective `UserRepo` methods.
    * Implement input validation (e.g., non-empty user name).
    * Handle errors, mapping repository errors to appropriate `connect.NewError` instances.
    * Convert between database model types and Protobuf message types.
    * Continue implementation until all tests in `user_test.go` pass.
9. **Register Service (`cmd/serve.go`):** Instantiate the `UserRepo`, instantiate the `UserService` (injecting the repo), and register the `UserService` handler with the Connect server.
10. **Run Tests:** Run `make test`. Ensure all tests in the project, including the new user tests, pass.
11. **Run Linter:** Run `make lint`. Fix all reported issues.
12. **Commit:** Commit changes with a descriptive message.

## Key Considerations

* Implement robust input validation for user data.
* Ensure proper error handling and mapping to Connect RPC errors.
* Implement pagination correctly for the `ListUsers` RPC.
* Follow the TDD cycle strictly: write failing test -> implement minimum code to pass -> refactor.
* Maintain consistency with the established project structure and coding standards.

## Implementation Order

1. Define Protobuf messages and service for Users.
2. Define SQL queries for User CRUD operations.
3. Run code generation (`make generate-all`).
4. Write failing service tests in `user_test.go`.
5. Implement `UserRepo` methods.
6. Implement `UserService` methods until tests pass.
7. Register `UserService` in `cmd/serve.go`.
8. Run `make test` and verify all tests pass.
9. Run `make lint` and fix issues.
10. Commit changes.
