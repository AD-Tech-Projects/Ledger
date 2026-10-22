-- Create "account_holders" table
CREATE TABLE "account_holders" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "unique_identifier" text NOT NULL,
  "holder_type" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_account_holders_deleted_at" to table: "account_holders"
CREATE INDEX "idx_account_holders_deleted_at" ON "account_holders" ("deleted_at");
-- Create index "idx_account_holders_unique_identifier" to table: "account_holders"
CREATE UNIQUE INDEX "idx_account_holders_unique_identifier" ON "account_holders" ("unique_identifier");
-- Create "transactions" table
CREATE TABLE "transactions" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "amount" bigint NOT NULL,
  "status" text NOT NULL,
  "transaction_type" text NOT NULL,
  "source_account_id" uuid NOT NULL,
  "destination_account_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_transactions_deleted_at" to table: "transactions"
CREATE INDEX "idx_transactions_deleted_at" ON "transactions" ("deleted_at");
-- Create index "idx_transactions_destination_account_id" to table: "transactions"
CREATE INDEX "idx_transactions_destination_account_id" ON "transactions" ("destination_account_id");
-- Create index "idx_transactions_source_account_id" to table: "transactions"
CREATE INDEX "idx_transactions_source_account_id" ON "transactions" ("source_account_id");
-- Create "accounts" table
CREATE TABLE "accounts" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "status" text NOT NULL,
  "account_type" text NOT NULL,
  "currency" text NOT NULL,
  "free_balance" bigint NOT NULL,
  "reserved_balance" bigint NOT NULL,
  "hold_balance" bigint NOT NULL,
  "holder_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_accounts_holder" FOREIGN KEY ("holder_id") REFERENCES "account_holders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_accounts_deleted_at" to table: "accounts"
CREATE INDEX "idx_accounts_deleted_at" ON "accounts" ("deleted_at");
-- Create index "idx_accounts_holder_id" to table: "accounts"
CREATE INDEX "idx_accounts_holder_id" ON "accounts" ("holder_id");
