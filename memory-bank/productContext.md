# Product Context: Expense Manager

## Purpose

To provide a technically sound backend for financial applications requiring accuracy via double-entry bookkeeping (DEB), serving individuals or small teams (**Users**) needing reliable financial tracking and potentially managing shared accounts.

## Problems Solved

1. **Financial Data Integrity:** Guarantees balanced transactions through DEB enforcement.
2. **Standard Accounting Foundation:** Enables accurate financial reporting.
3. **Structured Classification:** Clearly separates A/L/E accounts from I/E classifications (`categories`).
4. **Clear Account Association:** Manages potentially shared access to financial accounts by multiple **Users**.
5. **Reconciliation Aid:** Helps identify cash settlement needs.

## How It Should Work (Backend Logic - Final Architecture & Tools)

1. **DEB Core:** `transactions` capture events; `ledger_entries` record balanced debits/credits against `accounts` (A/L/E only).
2. **Account Structure:** `accounts` table holds A/L/E types. `name` is unique ID. Context from linked `instrument`, `institution`, and associated **`users`** (via `account_users`).
3. **Income/Expense Handling:** Balanced against a designated Equity account. `ledger_entries` affecting Equity link to a `categories.id`.
4. **Direct Service-Repo Interaction:** Service layer (`internal/rpc/services`) calls data access layer (`internal/repo`) implementations directly.
5. **Migrations (Atlas `schema.hcl` Diffing):** Schema managed via `db/schema.hcl`. Atlas generates versioned `UP` migration scripts.
6. **TDD Focus:** Service-level tests (using standard `testing`/`cmp`) drive implementation. Linting enforced.
7. **Configuration:** Loaded from environment variables via `carlosar/env`.

## User Experience Goals (Guiding Backend and Frontend Design)

The application (backend API and frontend UI) should:

1. **Simplify DEB:** The frontend UI should abstract complex accounting mechanics, presenting a user-friendly interface for financial operations, while the backend API enforces DEB rules.
2. **Ensure Integrity:** Backend validations prevent errors, and the frontend should guide users to provide valid input.
3. **Provide Clarity:** Structured data from the backend supports informative display in the frontend UI.
4. **Support Collaboration:** The **User** association model enables correct data presentation and potential filtering based on the logged-in user in the frontend.
5. **Provide a Responsive and Intuitive Interface:** The frontend should be easy to navigate and use on various devices.
