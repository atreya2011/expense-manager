-- Create "account_types" table
CREATE TABLE `account_types` (`id` text NULL, `name` text NOT NULL, `code` text NOT NULL, `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), CHECK (code IN ('A', 'L', 'E')));
-- Create index "account_types_name" to table: "account_types"
CREATE UNIQUE INDEX `account_types_name` ON `account_types` (`name`);
-- Create index "account_types_code" to table: "account_types"
CREATE UNIQUE INDEX `account_types_code` ON `account_types` (`code`);
-- Create "users" table
CREATE TABLE `users` (`id` text NULL, `name` text NOT NULL, `email` text NOT NULL, `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`));
-- Create index "users_email" to table: "users"
CREATE UNIQUE INDEX `users_email` ON `users` (`email`);
-- Create "instruments" table
CREATE TABLE `instruments` (`id` text NULL, `name` text NOT NULL, `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`));
-- Create index "instruments_name" to table: "instruments"
CREATE UNIQUE INDEX `instruments_name` ON `instruments` (`name`);
-- Create "currencies" table
CREATE TABLE `currencies` (`id` text NULL, `code` text NOT NULL, `name` text NOT NULL, `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`));
-- Create index "currencies_code" to table: "currencies"
CREATE UNIQUE INDEX `currencies_code` ON `currencies` (`code`);
-- Create "institutions" table
CREATE TABLE `institutions` (`id` text NULL, `name` text NOT NULL, `type` text NOT NULL, `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`));
-- Create index "institutions_name" to table: "institutions"
CREATE UNIQUE INDEX `institutions_name` ON `institutions` (`name`);
-- Create "accounts" table
CREATE TABLE `accounts` (`id` text NULL, `name` text NOT NULL, `description` text NULL, `account_type_id` text NOT NULL, `instrument_id` text NULL, `institution_id` text NULL, `currency_id` text NULL, `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `1` FOREIGN KEY (`institution_id`) REFERENCES `institutions` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `2` FOREIGN KEY (`instrument_id`) REFERENCES `instruments` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `3` FOREIGN KEY (`account_type_id`) REFERENCES `account_types` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "accounts_name" to table: "accounts"
CREATE UNIQUE INDEX `accounts_name` ON `accounts` (`name`);
-- Create "account_users" table
CREATE TABLE `account_users` (`account_id` text NOT NULL, `user_id` text NOT NULL, `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`account_id`, `user_id`), CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "categories" table
CREATE TABLE `categories` (`id` text NULL, `parent_id` text NULL, `name` text NOT NULL, `description` text NULL, `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`parent_id`) REFERENCES `categories` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "categories_name" to table: "categories"
CREATE UNIQUE INDEX `categories_name` ON `categories` (`name`);
-- Create "transactions" table
CREATE TABLE `transactions` (`id` text NULL, `date` timestamp NOT NULL, `description` text NOT NULL, `notes` text NULL, `category_id` text NULL, `instrument_id` text NULL, `allocation_tag` text NULL, `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`instrument_id`) REFERENCES `instruments` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `1` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create "ledger_entries" table
CREATE TABLE `ledger_entries` (`id` text NULL, `transaction_id` text NOT NULL, `account_id` text NOT NULL, `category_id` text NULL, `memo` text NOT NULL, `debit` integer NOT NULL DEFAULT 0, `credit` integer NOT NULL DEFAULT 0, `currency_id` text NULL, `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`currency_id`) REFERENCES `currencies` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `1` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `2` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT `3` FOREIGN KEY (`transaction_id`) REFERENCES `transactions` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CHECK (
    debit >= 0
    AND credit >= 0
    AND (
      debit = 0
      OR credit = 0
    )
  ));
