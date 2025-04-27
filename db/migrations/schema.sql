-- Account Types (Asset, Liability, Equity)
CREATE TABLE account_types (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  code TEXT NOT NULL CHECK (code IN ('A', 'L', 'E')),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (name),
  UNIQUE (code)
);

-- Users
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (email)
);

-- Instruments (Cash, Bank Account, Credit Card, etc.)
CREATE TABLE instruments (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (name)
);

-- Currencies
CREATE TABLE currencies (
  id TEXT PRIMARY KEY,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (code)
);

-- Institutions (Banks, Credit Card Companies, etc.)
CREATE TABLE institutions (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (name)
);

-- Accounts
CREATE TABLE accounts (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  account_type_id TEXT NOT NULL,
  instrument_id TEXT,
  institution_id TEXT,
  currency_id TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (name),
  FOREIGN KEY (account_type_id) REFERENCES account_types (id),
  FOREIGN KEY (instrument_id) REFERENCES instruments (id),
  FOREIGN KEY (institution_id) REFERENCES institutions (id),
  FOREIGN KEY (currency_id) REFERENCES currencies (id)
);

-- Account Users (Many-to-Many)
CREATE TABLE account_users (
  account_id TEXT NOT NULL,
  user_id TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (account_id, user_id),
  FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Categories (Income/Expense types, optional parent_id for hierarchy)
CREATE TABLE categories (
  id TEXT PRIMARY KEY,
  parent_id TEXT,
  name TEXT NOT NULL,
  description TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (name),
  FOREIGN KEY (parent_id) REFERENCES categories (id)
);

-- Transactions (Journal)
CREATE TABLE transactions (
  id TEXT PRIMARY KEY,
  date TIMESTAMP NOT NULL,
  description TEXT NOT NULL,
  notes TEXT,
  category_id TEXT,
  instrument_id TEXT,
  allocation_tag TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (category_id) REFERENCES categories (id),
  FOREIGN KEY (instrument_id) REFERENCES instruments (id)
);

-- Ledger Entries
CREATE TABLE ledger_entries (
  id TEXT PRIMARY KEY,
  transaction_id TEXT NOT NULL,
  account_id TEXT NOT NULL,
  category_id TEXT,
  memo TEXT NOT NULL,
  debit INTEGER NOT NULL DEFAULT 0,
  credit INTEGER NOT NULL DEFAULT 0,
  currency_id TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (transaction_id) REFERENCES transactions (id) ON DELETE CASCADE,
  FOREIGN KEY (account_id) REFERENCES accounts (id),
  FOREIGN KEY (category_id) REFERENCES categories (id),
  FOREIGN KEY (currency_id) REFERENCES currencies (id),
  CHECK (
    debit >= 0
    AND credit >= 0
    AND (
      debit = 0
      OR credit = 0
    )
  )
);
