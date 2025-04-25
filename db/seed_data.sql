-- Insert account types
INSERT INTO
  account_types (id, name, code)
VALUES
  ('at_asset', 'Asset', 'A'),
  ('at_liability', 'Liability', 'L'),
  ('at_equity', 'Equity', 'E');

-- Insert default currency (JPY)
INSERT INTO
  currencies (id, code, name)
VALUES
  ('cur_jpy', 'JPY', 'Japanese Yen');

-- Insert basic instruments
INSERT INTO
  instruments (id, name)
VALUES
  ('inst_cash', 'Cash'),
  ('inst_bank', 'Bank Account'),
  ('inst_credit', 'Credit Card');

-- Insert essential Equity account for tracking earnings
INSERT INTO
  accounts (id, name, description, account_type_id)
VALUES
  (
    'acc_earnings',
    'Current Year Earnings',
    'Account for tracking current year earnings',
    'at_equity'
  );
