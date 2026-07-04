# Finance Database

`finance-service` can connect to PostgreSQL and store transactions through:

```text
finance.v1.FinanceService/CreateTransaction
```

Required environment variables:

```bash
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=<local_database_user>
DATABASE_PASSWORD=<local_database_password>
DATABASE_NAME=my_skripsi_gwejh
DATABASE_SSLMODE=disable
```

When the service starts with `DATABASE_USER` and `DATABASE_NAME`, it runs the embedded schema and seed SQL.

## Local Run

```bash
cd finance-service
DATABASE_HOST=localhost \
DATABASE_PORT=5432 \
DATABASE_USER=<local_database_user> \
DATABASE_PASSWORD=<local_database_password> \
DATABASE_NAME=my_skripsi_gwejh \
DATABASE_SSLMODE=disable \
GRPC_PORT=50053 \
go run ./cmd/finance-service
```

## Get Seed Data

```sql
SELECT account_id FROM account WHERE whatsapp_number = '6281234567890';
SELECT wallet_id FROM wallet WHERE is_primary = true LIMIT 1;
SELECT category_id FROM category WHERE name = 'makanan' LIMIT 1;
```

## Create Transaction

```bash
grpcurl -plaintext \
  -d '{
    "account_id": "ACCOUNT_UUID_DUMMY",
    "wallet_id": 1,
    "category_id": 1,
    "type": "expense",
    "amount": 20000,
    "name": "bayar kopi",
    "is_ai_generated": true,
    "report_date": "2026-07-04T10:00:00Z"
  }' \
  localhost:50053 \
  finance.v1.FinanceService/CreateTransaction
```

## Verify Data

```sql
SELECT
  transaction_id,
  account_id,
  wallet_id,
  category_id,
  type,
  amount,
  name,
  is_ai_generated,
  report_date,
  created_at
FROM "transaction"
ORDER BY transaction_id DESC
LIMIT 5;
```
