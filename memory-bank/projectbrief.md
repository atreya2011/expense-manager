# Project Brief: Expense Manager

## Project Overview

Expense Manager is a backend API service providing accurate financial tracking using double-entry bookkeeping. Built with Go, ConnectRPC, SQLite, and modern tooling (Buf, sqlc, Atlas, Air), it emphasizes data integrity, a TDD workflow, and mandatory code linting.

## Core Requirements

1. **Double-Entry Bookkeeping Core:** Enforce balanced debit/credit ledger entries.
2. **Account Management:** Support Asset, Liability, Equity accounts. Manage association with **Users** (M2M) and link to institutions/instruments. Use `accounts.name` as UNIQUE key.
3. **Transaction Journaling:** Record event metadata.
4. **Ledger Management:** Detail financial impact on A/L/E accounts via `ledger_entries`.
5. **Income/Expense Classification:** Use a separate `categories` table (optional hierarchy) linked to ledger entries affecting Equity.
6. **User Tracking:** Support multiple Users associated with accounts.
7. **Basic Reporting:** API endpoints for account balance and category summary calculations.
8. **Reconciliation Support:** API/logic to identify cash reconciliation needs.

## Project Scope (Initial - Backend API)

* Develop Go/ConnectRPC backend service.
* Implement DB schema using SQLite, managed via Atlas (`schema.hcl` diffing, versioned UP migrations).
* Provide API endpoints for managing master data (**Users**, Institutions, Instruments, Currencies, AccountTypes, Categories) and core entities (Accounts [A/L/E], Transactions, Ledger Entries).
* Enforce DEB rules.
* Implement core reporting and reconciliation ID endpoints.
* Focus on JPY currency.
* Utilize Buf, sqlc, Air.
* **Mandatory `golangci-lint`**.
* **Employ Test-Driven Development (TDD)** using standard `testing` and `google/go-cmp/cmp`.
* **Configuration via environment variables** (`carlosar/env`).

## Out of Scope (Initial)

* Budgeting, Investments, Automated Reconciliation.
* Advanced multi-currency features.
* User Authentication/Authorization (beyond basic User records).
* Database down migrations.

## Future Considerations

* Frontend applications (planning complete, implementation to be interleaved with backend).

* Features listed as Out of Scope.
* AuthN/AuthZ implementation.
* Performance tuning.
