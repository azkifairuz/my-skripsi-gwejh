CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transaction_type') THEN
    CREATE TYPE transaction_type AS ENUM ('income', 'expense');
  END IF;
END
$$;

CREATE TABLE IF NOT EXISTS account (
  account_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR,
  password VARCHAR,
  whatsapp_number VARCHAR UNIQUE,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wallet (
  wallet_id SERIAL PRIMARY KEY,
  account_id UUID REFERENCES account(account_id),
  name VARCHAR,
  balance DECIMAL DEFAULT 0,
  is_primary BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS category (
  category_id SERIAL PRIMARY KEY,
  account_id UUID REFERENCES account(account_id),
  name VARCHAR,
  icon VARCHAR,
  UNIQUE (account_id, name)
);

CREATE UNIQUE INDEX IF NOT EXISTS category_global_name_idx
  ON category (name)
  WHERE account_id IS NULL;

CREATE TABLE IF NOT EXISTS budget (
  budget_id SERIAL PRIMARY KEY,
  account_id UUID REFERENCES account(account_id),
  category_id INTEGER REFERENCES category(category_id),
  name VARCHAR,
  amount DECIMAL,
  period VARCHAR
);

CREATE TABLE IF NOT EXISTS "transaction" (
  transaction_id SERIAL PRIMARY KEY,
  account_id UUID REFERENCES account(account_id),
  wallet_id INTEGER REFERENCES wallet(wallet_id),
  category_id INTEGER REFERENCES category(category_id),
  budget_id INTEGER REFERENCES budget(budget_id),
  type transaction_type,
  amount DECIMAL,
  name VARCHAR,
  is_ai_generated BOOLEAN DEFAULT false,
  receipt_image_url VARCHAR,
  report_date TIMESTAMP,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP
);
