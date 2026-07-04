WITH dummy_account AS (
  INSERT INTO account (username, password, whatsapp_number)
  VALUES ('dummy_user', 'dummy_password', '6281234567890')
  ON CONFLICT (whatsapp_number) DO UPDATE
    SET username = EXCLUDED.username
  RETURNING account_id
),
existing_dummy_account AS (
  SELECT account_id FROM account WHERE whatsapp_number = '6281234567890'
),
selected_account AS (
  SELECT account_id FROM dummy_account
  UNION
  SELECT account_id FROM existing_dummy_account
  LIMIT 1
)
INSERT INTO wallet (account_id, name, balance, is_primary)
SELECT account_id, 'Wallet Utama', 0, true
FROM selected_account
WHERE NOT EXISTS (
  SELECT 1
  FROM wallet
  WHERE account_id = (SELECT account_id FROM selected_account)
    AND is_primary = true
);

INSERT INTO category (account_id, name, icon)
VALUES
  (NULL, 'makanan', NULL),
  (NULL, 'transportasi', NULL),
  (NULL, 'belanja', NULL),
  (NULL, 'tagihan', NULL),
  (NULL, 'pemasukan', NULL),
  (NULL, 'kesehatan', NULL),
  (NULL, 'pendidikan', NULL),
  (NULL, 'hiburan', NULL),
  (NULL, 'lainnya', NULL)
ON CONFLICT (name) WHERE account_id IS NULL DO NOTHING;
