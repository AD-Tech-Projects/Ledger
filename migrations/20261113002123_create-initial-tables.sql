-- Create "owners" table
CREATE TABLE "owners" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "unique_identifier" text NOT NULL,
  "owner_type" text NOT NULL,
  "country_code" character varying(2) NOT NULL DEFAULT '',
  PRIMARY KEY ("id")
);
-- Create index "idx_owners_deleted_at" to table: "owners"
CREATE INDEX "idx_owners_deleted_at" ON "owners" ("deleted_at");
-- Create index "idx_owners_unique_identifier" to table: "owners"
CREATE UNIQUE INDEX "idx_owners_unique_identifier" ON "owners" ("unique_identifier");
-- Create "external_accounts" table
CREATE TABLE "external_accounts" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "owner_id" uuid NOT NULL,
  "kind" text NOT NULL,
  "country_code" character varying(2) NULL,
  "currency" text NOT NULL,
  "institution_name" text NULL,
  "branch_code" text NULL,
  "account_number" text NULL,
  "routing_code" text NULL,
  "iban" text NULL,
  "bic" text NULL,
  "crypto_network" text NULL,
  "crypto_address" text NULL,
  "crypto_memo" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_external_accounts_owner" FOREIGN KEY ("owner_id") REFERENCES "owners" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_external_accounts_deleted_at" to table: "external_accounts"
CREATE INDEX "idx_external_accounts_deleted_at" ON "external_accounts" ("deleted_at");
-- Create index "idx_external_accounts_owner_currency" to table: "external_accounts"
CREATE INDEX "idx_external_accounts_owner_currency" ON "external_accounts" ("owner_id", "currency");
-- Create "account_identifiers" table
CREATE TABLE "account_identifiers" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "account_id" uuid NOT NULL,
  "network" text NOT NULL,
  "identifier_type" text NOT NULL,
  "value" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_external_accounts_identifiers" FOREIGN KEY ("account_id") REFERENCES "external_accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_account_identifiers_account_id" to table: "account_identifiers"
CREATE INDEX "idx_account_identifiers_account_id" ON "account_identifiers" ("account_id");
-- Create index "idx_account_identifiers_deleted_at" to table: "account_identifiers"
CREATE INDEX "idx_account_identifiers_deleted_at" ON "account_identifiers" ("deleted_at");
-- Create "assets" table
CREATE TABLE "assets" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "name" text NOT NULL,
  "type" text NOT NULL,
  "status" text NOT NULL,
  "code" text NOT NULL,
  "decimal_places" smallint NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_assets_code" to table: "assets"
CREATE UNIQUE INDEX "idx_assets_code" ON "assets" ("code");
-- Create index "idx_assets_deleted_at" to table: "assets"
CREATE INDEX "idx_assets_deleted_at" ON "assets" ("deleted_at");
-- Create "asset_prices" table
CREATE TABLE "asset_prices" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "price" bigint NOT NULL,
  "recorded_at" timestamptz NOT NULL,
  "asset_id" uuid NOT NULL,
  "quote_asset_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_asset_prices_asset" FOREIGN KEY ("asset_id") REFERENCES "assets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_asset_prices_quote_asset" FOREIGN KEY ("quote_asset_id") REFERENCES "assets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_asset_prices_asset_id" to table: "asset_prices"
CREATE INDEX "idx_asset_prices_asset_id" ON "asset_prices" ("asset_id");
-- Create index "idx_asset_prices_deleted_at" to table: "asset_prices"
CREATE INDEX "idx_asset_prices_deleted_at" ON "asset_prices" ("deleted_at");
-- Create index "idx_asset_prices_quote_asset_id" to table: "asset_prices"
CREATE INDEX "idx_asset_prices_quote_asset_id" ON "asset_prices" ("quote_asset_id");
-- Create index "idx_asset_prices_recorded_at" to table: "asset_prices"
CREATE INDEX "idx_asset_prices_recorded_at" ON "asset_prices" ("recorded_at");
-- Create "ledger_accounts" table
CREATE TABLE "ledger_accounts" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "type" text NOT NULL,
  "status" text NOT NULL,
  "normal_balance_direction" text NOT NULL,
  "name" text NOT NULL,
  "code" text NOT NULL,
  "asset_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_ledger_accounts_asset" FOREIGN KEY ("asset_id") REFERENCES "assets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_ledger_accounts_asset_id" to table: "ledger_accounts"
CREATE INDEX "idx_ledger_accounts_asset_id" ON "ledger_accounts" ("asset_id");
-- Create index "idx_ledger_accounts_code" to table: "ledger_accounts"
CREATE UNIQUE INDEX "idx_ledger_accounts_code" ON "ledger_accounts" ("code");
-- Create index "idx_ledger_accounts_deleted_at" to table: "ledger_accounts"
CREATE INDEX "idx_ledger_accounts_deleted_at" ON "ledger_accounts" ("deleted_at");
-- Create "wallets" table
CREATE TABLE "wallets" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "type" text NOT NULL,
  "status" text NOT NULL,
  "version" bigint NOT NULL,
  "free_balance" bigint NOT NULL,
  "reserved_balance" bigint NOT NULL,
  "hold_balance" bigint NOT NULL,
  "owner_id" uuid NOT NULL,
  "asset_id" uuid NOT NULL,
  "account_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_wallets_account" FOREIGN KEY ("account_id") REFERENCES "external_accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_wallets_asset" FOREIGN KEY ("asset_id") REFERENCES "assets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_wallets_owner" FOREIGN KEY ("owner_id") REFERENCES "owners" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_wallets_free_balance" CHECK (free_balance >= 0),
  CONSTRAINT "chk_wallets_hold_balance" CHECK (hold_balance >= 0),
  CONSTRAINT "chk_wallets_reserved_balance" CHECK (reserved_balance >= 0)
);
-- Create index "idx_wallets_account_id" to table: "wallets"
CREATE INDEX "idx_wallets_account_id" ON "wallets" ("account_id");
-- Create index "idx_wallets_asset_id" to table: "wallets"
CREATE INDEX "idx_wallets_asset_id" ON "wallets" ("asset_id");
-- Create index "idx_wallets_deleted_at" to table: "wallets"
CREATE INDEX "idx_wallets_deleted_at" ON "wallets" ("deleted_at");
-- Create index "idx_wallets_owner_id" to table: "wallets"
CREATE INDEX "idx_wallets_owner_id" ON "wallets" ("owner_id");
-- Create "transactions" table
CREATE TABLE "transactions" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "amount" bigint NOT NULL,
  "status" text NOT NULL,
  "transaction_type" text NOT NULL,
  "asset_id" uuid NOT NULL,
  "source_wallet_id" uuid NULL,
  "destination_wallet_id" uuid NULL,
  "reference_id" text NULL,
  "reversal_origin_id" uuid NULL,
  "description" text NOT NULL,
  "completed_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_transactions_asset" FOREIGN KEY ("asset_id") REFERENCES "assets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_transactions_destination_wallet" FOREIGN KEY ("destination_wallet_id") REFERENCES "wallets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_transactions_reversal_origin" FOREIGN KEY ("reversal_origin_id") REFERENCES "transactions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_transactions_source_wallet" FOREIGN KEY ("source_wallet_id") REFERENCES "wallets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_transactions_amount" CHECK (amount > 0)
);
-- Create index "idx_transactions_asset_id" to table: "transactions"
CREATE INDEX "idx_transactions_asset_id" ON "transactions" ("asset_id");
-- Create index "idx_transactions_deleted_at" to table: "transactions"
CREATE INDEX "idx_transactions_deleted_at" ON "transactions" ("deleted_at");
-- Create index "idx_transactions_destination_wallet_id" to table: "transactions"
CREATE INDEX "idx_transactions_destination_wallet_id" ON "transactions" ("destination_wallet_id");
-- Create index "idx_transactions_reference_id" to table: "transactions"
CREATE INDEX "idx_transactions_reference_id" ON "transactions" ("reference_id");
-- Create index "idx_transactions_reversal_origin_id" to table: "transactions"
CREATE INDEX "idx_transactions_reversal_origin_id" ON "transactions" ("reversal_origin_id");
-- Create index "idx_transactions_source_wallet_id" to table: "transactions"
CREATE INDEX "idx_transactions_source_wallet_id" ON "transactions" ("source_wallet_id");
-- Create "ledger_transactions" table
CREATE TABLE "ledger_transactions" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "transaction_id" uuid NOT NULL,
  "status" text NOT NULL,
  "description" text NOT NULL,
  "posted_at" timestamptz NULL,
  "reversal_origin_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_ledger_transactions_reversal_origin" FOREIGN KEY ("reversal_origin_id") REFERENCES "ledger_transactions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_ledger_transactions_transaction" FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_ledger_transactions_deleted_at" to table: "ledger_transactions"
CREATE INDEX "idx_ledger_transactions_deleted_at" ON "ledger_transactions" ("deleted_at");
-- Create index "idx_ledger_transactions_reversal_origin_id" to table: "ledger_transactions"
CREATE INDEX "idx_ledger_transactions_reversal_origin_id" ON "ledger_transactions" ("reversal_origin_id");
-- Create index "idx_ledger_transactions_transaction_id" to table: "ledger_transactions"
CREATE UNIQUE INDEX "idx_ledger_transactions_transaction_id" ON "ledger_transactions" ("transaction_id");
-- Create "ledger_entries" table
CREATE TABLE "ledger_entries" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "ledger_transaction_id" uuid NOT NULL,
  "ledger_account_id" uuid NOT NULL,
  "asset_id" uuid NOT NULL,
  "amount" bigint NOT NULL,
  "direction" text NOT NULL,
  "sequence" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_ledger_entries_asset" FOREIGN KEY ("asset_id") REFERENCES "assets" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_ledger_entries_ledger_account" FOREIGN KEY ("ledger_account_id") REFERENCES "ledger_accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_ledger_entries_ledger_transaction" FOREIGN KEY ("ledger_transaction_id") REFERENCES "ledger_transactions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_ledger_entries_amount" CHECK (amount > 0),
  CONSTRAINT "chk_ledger_entries_direction" CHECK (direction = ANY (ARRAY['DEBIT'::text, 'CREDIT'::text]))
);
-- Create index "idx_ledger_entries_asset_id" to table: "ledger_entries"
CREATE INDEX "idx_ledger_entries_asset_id" ON "ledger_entries" ("asset_id");
-- Create index "idx_ledger_entries_deleted_at" to table: "ledger_entries"
CREATE INDEX "idx_ledger_entries_deleted_at" ON "ledger_entries" ("deleted_at");
-- Create index "idx_ledger_entries_ledger_account_id" to table: "ledger_entries"
CREATE INDEX "idx_ledger_entries_ledger_account_id" ON "ledger_entries" ("ledger_account_id");
-- Create index "idx_ledger_entries_transaction_sequence" to table: "ledger_entries"
CREATE UNIQUE INDEX "idx_ledger_entries_transaction_sequence" ON "ledger_entries" ("ledger_transaction_id", "sequence");
