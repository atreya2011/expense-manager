# Phase 2 Implementation Plan: Instruments Service

**Instructions for LLM:**
1. After completing each numbered step in the "Implementation Steps" section, wait for explicit approval from the user before proceeding to the next step.
2. After all steps in this file are completed and approved, update the memory bank to reflect the completed work for this phase.
3. Finally, use the `attempt_completion` tool to summarize the work done for this phase.

## Overview

This document outlines the plan for implementing the Instruments Service as part of Phase 2 of the Expense Manager project. According to the implementation plan, this phase should focus on implementing read/list functionality for `instruments`, with full CRUD if needed.

## Implementation Steps

### 1. Define Protobuf Messages and Service

Create a new proto file `proto/expenses/v1/instrument.proto` with the following components:

- **Messages:**
  - `Instrument`: Core message with fields matching the database schema (id, name, created_at, updated_at)
  - Request/Response messages for operations:
    - `GetInstrumentRequest`/`GetInstrumentResponse`
    - `ListInstrumentsRequest`/`ListInstrumentsResponse`
    - `CreateInstrumentRequest`/`CreateInstrumentResponse`
    - `UpdateInstrumentRequest`/`UpdateInstrumentResponse`
    - `DeleteInstrumentRequest`/`DeleteInstrumentResponse`

- **Service Definition:**
  - `InstrumentService` with RPCs:
    - `GetInstrument`
    - `ListInstruments`
    - `CreateInstrument`
    - `UpdateInstrument`
    - `DeleteInstrument`

Include appropriate imports for common types like Pagination and Timestamp.

### 2. Define SQL Queries

Create a new SQL file `db/queries/instrument.sql` with the following queries:

- `GetInstrument`: Query to retrieve a single instrument by ID
- `ListInstruments`: Query to retrieve a paginated list of instruments
- `CreateInstrument`: Query to create a new instrument
- `UpdateInstrument`: Query to update an existing instrument
- `DeleteInstrument`: Query to delete an instrument by ID

Each query should include appropriate sqlc annotations.

### 3. Generate Code

After defining the proto and SQL files, run `make generate-all` to generate:
- Proto stubs in `internal/rpc/gen/expenses/v1/`
- SQLc generated code in `internal/repo/gen/`

### 4. Write Tests

Create a test file `internal/rpc/services/instrument_test.go` with comprehensive test cases:

- TestCreateInstrument: Test instrument creation
- TestGetInstrument: Test retrieving an instrument by ID
- TestListInstruments: Test listing instruments with pagination
- TestUpdateInstrument: Test updating an instrument
- TestDeleteInstrument: Test deleting an instrument

Ensure tests cover both successful operations and error cases.

### 5. Implement Repository

Create a repository file `internal/repo/instrument.go` with the following methods:

- `NewInstrumentRepo`: Constructor function
- `CreateInstrument`: Create a new instrument
- `GetInstrument`: Retrieve an instrument by ID
- `ListInstruments`: List instruments with pagination
- `UpdateInstrument`: Update an existing instrument
- `DeleteInstrument`: Delete an instrument by ID

Follow the pattern established in the existing UserRepo implementation.

### 6. Implement Service

Create a service file `internal/rpc/services/instrument.go` implementing the InstrumentService interface:

- `NewInstrumentService`: Constructor function
- `CreateInstrument`: Create a new instrument
- `GetInstrument`: Retrieve an instrument by ID
- `ListInstruments`: List instruments with pagination
- `UpdateInstrument`: Update an existing instrument
- `DeleteInstrument`: Delete an instrument by ID

Include proper validation, error handling, and conversion between domain and proto objects.

### 7. Register Service

Update `cmd/serve.go` to:
- Initialize the InstrumentRepo
- Create an instance of InstrumentService
- Register the service with the Connect RPC server

### 8. Test and Lint

Run `make test` to ensure all tests pass and `make lint` to verify code quality.

## Key Considerations

1. **Data Validation**: Validate input parameters, particularly ensuring instrument names are unique and not empty.

2. **Error Handling**: Follow the established pattern of mapping database errors to appropriate Connect RPC error codes.

3. **Pagination**: Implement pagination for ListInstruments similar to the User service.

4. **Consistency**: Maintain consistency with existing patterns in the User service implementation.

5. **TDD Approach**: Follow the Test-Driven Development workflow by writing tests first, implementing the minimal required code to make tests pass, and then refactoring as needed.

## Implementation Order

1. Define proto and SQL files
2. Generate code
3. Write failing tests
4. Implement repository
5. Implement service
6. Update server registration
7. Verify tests pass
8. Fix any linting issues

This implementation follows the Standard Service Implementation Workflow defined in the overall implementation plan and maintains consistency with the patterns established in the User service.
