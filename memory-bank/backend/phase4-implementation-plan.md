# Phase 4 Implementation Plan: Institutions Service

**Instructions for LLM:**
1. After completing each numbered step in the "Implementation Steps" section, wait for explicit approval from the user before proceeding to the next step.
2. After all steps in this file are completed and approved, update the memory bank to reflect the completed work for this phase.
3. Finally, use the `attempt_completion` tool to summarize the work done for this phase.

## Overview

This phase focuses on implementing the full CRUD (Create, Read, Update, Delete) functionality for the `institutions` table.

## Implementation Steps

This phase follows the **Standard Service Implementation Workflow (TDD)**:

1. **Define API (`proto/expenses/v1/expenses.proto`):** Define `Institution` message, and request/response messages for `CreateInstitution`, `GetInstitution`, `ListInstitutions`, `UpdateInstitution`, `DeleteInstitution` RPCs. Add `service InstitutionService {}` to `expenses.proto`.
2. **Define SQL (`db/queries/institution.sql`):** Write `sqlc`-annotated queries for `CreateInstitution`, `GetInstitution`, `ListInstitutions`, `UpdateInstitution`, `DeleteInstitution`.
3. **Define Schema (`db/schema.sql`):** (Already done in Phase 0, includes `institutions` table). No changes needed unless schema adjustments are required.
4. **Generate Code:** Run `make generate-all`. Verify generated Go code.
5. **Create/Apply Migration:** (Only if schema changed in step 3). Run `make migrate-new NAME=...` followed by `make migrate`.
6. **Write Failing Service Test (`internal/rpc/services/institution_test.go`):**
    * Create the test file.
    * Implement `TestMain` for test database setup/teardown.
    * Write table-driven tests for all CRUD operations.
    * Instantiate concrete repo and service.
    * Use `google/go-cmp/cmp` for assertions.
    * Run tests and confirm they initially fail.
7. **Implement Repository Logic (`internal/repo/institution.go`):**
    * Create the `InstitutionRepo` struct.
    * Implement CRUD methods calling `sqlc`-generated functions.
8. **Implement Service Logic (`internal/rpc/services/institution.go`):**
    * Create the `InstitutionService` struct.
    * Implement the `InstitutionService` interface methods.
    * Call `InstitutionRepo` methods.
    * Implement input validation (e.g., non-empty name).
    * Handle errors.
    * Convert between types.
    * Continue until tests pass.
9. **Register Service (`cmd/serve.go`):** Instantiate and register the `InstitutionRepo` and `InstitutionService`.
10. **Run Tests:** Run `make test`. Ensure all tests pass.
11. **Run Linter:** Run `make lint`. Fix all reported issues.
12. **Commit:** Commit changes.

## Key Considerations

* Implement input validation for institution data.
* Ensure proper error handling, especially for cases like deleting an institution that is still linked to accounts.
* Implement pagination for `ListInstitutions`.

## Implementation Order

1. Define Protobuf messages and service for Institutions.
2. Define SQL queries for Institution CRUD operations.
3. Run code generation (`make generate-all`).
4. Write failing service tests in `institution_test.go`.
5. Implement `InstitutionRepo` methods.
6. Implement `InstitutionService` methods until tests pass.
7. Register `InstitutionService` in `cmd/serve.go`.
8. Run `make test` and verify all tests pass.
9. Run `make lint` and fix issues.
10. Commit changes.
