package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrAccountNotFound = errors.New("account not found")

type RegisterAccountParams struct {
	WhatsappNumber string
	Username       string
}

type Account struct {
	AccountID       string
	WhatsappNumber  string
	Username        string
	PrimaryWalletID int64
	Created         bool
	CreatedAt       time.Time
}

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) RegisterByWhatsappNumber(ctx context.Context, params RegisterAccountParams) (*Account, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	account, err := upsertAccount(ctx, tx, params)
	if err != nil {
		return nil, err
	}

	walletID, err := ensurePrimaryWallet(ctx, tx, account.AccountID)
	if err != nil {
		return nil, err
	}
	account.PrimaryWalletID = walletID

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return account, nil
}

func (r *AccountRepository) LoginByWhatsappNumber(ctx context.Context, whatsappNumber string) (*Account, error) {
	query := `
		SELECT
			a.account_id::text,
			a.whatsapp_number,
			COALESCE(a.username, ''),
			COALESCE(w.wallet_id, 0),
			a.created_at
		FROM account a
		LEFT JOIN wallet w
			ON w.account_id = a.account_id
			AND w.is_primary = true
		WHERE a.whatsapp_number = $1
		LIMIT 1
	`

	var account Account
	err := r.db.QueryRowContext(ctx, query, whatsappNumber).Scan(
		&account.AccountID,
		&account.WhatsappNumber,
		&account.Username,
		&account.PrimaryWalletID,
		&account.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAccountNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find account by whatsapp number: %w", err)
	}

	return &account, nil
}

func upsertAccount(ctx context.Context, tx *sql.Tx, params RegisterAccountParams) (*Account, error) {
	query := `
		WITH inserted AS (
			INSERT INTO account (username, whatsapp_number)
			VALUES (NULLIF($1, ''), $2)
			ON CONFLICT (whatsapp_number) DO NOTHING
			RETURNING account_id::text, whatsapp_number, COALESCE(username, ''), created_at, true AS created
		)
		SELECT account_id, whatsapp_number, username, created_at, created FROM inserted
		UNION ALL
		SELECT account_id::text, whatsapp_number, COALESCE(username, ''), created_at, false AS created
		FROM account
		WHERE whatsapp_number = $2
		  AND NOT EXISTS (SELECT 1 FROM inserted)
		LIMIT 1
	`

	var account Account
	err := tx.QueryRowContext(ctx, query, params.Username, params.WhatsappNumber).Scan(
		&account.AccountID,
		&account.WhatsappNumber,
		&account.Username,
		&account.CreatedAt,
		&account.Created,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert account: %w", err)
	}

	return &account, nil
}

func ensurePrimaryWallet(ctx context.Context, tx *sql.Tx, accountID string) (int64, error) {
	query := `
		WITH inserted AS (
			INSERT INTO wallet (account_id, name, balance, is_primary)
			SELECT $1::uuid, 'Wallet Utama', 0, true
			WHERE NOT EXISTS (
				SELECT 1 FROM wallet WHERE account_id = $1::uuid AND is_primary = true
			)
			RETURNING wallet_id
		)
		SELECT wallet_id FROM inserted
		UNION ALL
		SELECT wallet_id FROM wallet
		WHERE account_id = $1::uuid AND is_primary = true
		  AND NOT EXISTS (SELECT 1 FROM inserted)
		LIMIT 1
	`

	var walletID int64
	if err := tx.QueryRowContext(ctx, query, accountID).Scan(&walletID); err != nil {
		return 0, fmt.Errorf("ensure primary wallet: %w", err)
	}

	return walletID, nil
}
