package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CreateTransactionParams struct {
	AccountID       string
	WalletID        int64
	CategoryID      int64
	BudgetID        *int64
	Type            string
	Amount          float64
	Name            string
	IsAIGenerated   bool
	ReceiptImageURL string
	ReportDate      time.Time
}

type ListTransactionHistoryParams struct {
	AccountID string
	Limit     int32
	Offset    int32
	Type      string
	FromDate  *time.Time
	ToDate    *time.Time
}

type Transaction struct {
	TransactionID   int64
	AccountID       string
	WalletID        int64
	WalletName      string
	CategoryID      int64
	CategoryName    string
	Type            string
	Amount          float64
	Name            string
	IsAIGenerated   bool
	ReceiptImageURL string
	ReportDate      time.Time
	CreatedAt       time.Time
}

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, params CreateTransactionParams) (*Transaction, error) {
	query := `
		INSERT INTO "transaction" (
			account_id,
			wallet_id,
			category_id,
			budget_id,
			type,
			amount,
			name,
			is_ai_generated,
			receipt_image_url,
			report_date
		)
		VALUES ($1, $2, $3, $4, $5::transaction_type, $6, $7, $8, $9, $10)
		RETURNING
			transaction_id,
			account_id::text,
			wallet_id,
			category_id,
			type::text,
			amount::float8,
			name,
			report_date,
			created_at
	`

	var budgetID sql.NullInt64
	if params.BudgetID != nil {
		budgetID = sql.NullInt64{Int64: *params.BudgetID, Valid: true}
	}

	var transaction Transaction
	err := r.db.QueryRowContext(
		ctx,
		query,
		params.AccountID,
		params.WalletID,
		params.CategoryID,
		budgetID,
		params.Type,
		params.Amount,
		params.Name,
		params.IsAIGenerated,
		params.ReceiptImageURL,
		params.ReportDate,
	).Scan(
		&transaction.TransactionID,
		&transaction.AccountID,
		&transaction.WalletID,
		&transaction.CategoryID,
		&transaction.Type,
		&transaction.Amount,
		&transaction.Name,
		&transaction.ReportDate,
		&transaction.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert transaction: %w", err)
	}

	return &transaction, nil
}

func (r *TransactionRepository) ListHistory(ctx context.Context, params ListTransactionHistoryParams) ([]Transaction, error) {
	query := `
		SELECT
			t.transaction_id,
			t.account_id::text,
			t.wallet_id,
			COALESCE(w.name, ''),
			t.category_id,
			COALESCE(c.name, ''),
			t.type::text,
			t.amount::float8,
			t.name,
			COALESCE(t.is_ai_generated, false),
			COALESCE(t.receipt_image_url, ''),
			t.report_date,
			t.created_at
		FROM "transaction" t
		LEFT JOIN wallet w ON w.wallet_id = t.wallet_id
		LEFT JOIN category c ON c.category_id = t.category_id
		WHERE t.account_id = $1::uuid
		  AND ($2 = '' OR t.type = $2::transaction_type)
		  AND ($3::timestamp IS NULL OR t.report_date >= $3)
		  AND ($4::timestamp IS NULL OR t.report_date <= $4)
		ORDER BY t.report_date DESC, t.transaction_id DESC
		LIMIT $5 OFFSET $6
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		params.AccountID,
		params.Type,
		params.FromDate,
		params.ToDate,
		params.Limit,
		params.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list transaction history: %w", err)
	}
	defer rows.Close()

	transactions := make([]Transaction, 0)
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(
			&transaction.TransactionID,
			&transaction.AccountID,
			&transaction.WalletID,
			&transaction.WalletName,
			&transaction.CategoryID,
			&transaction.CategoryName,
			&transaction.Type,
			&transaction.Amount,
			&transaction.Name,
			&transaction.IsAIGenerated,
			&transaction.ReceiptImageURL,
			&transaction.ReportDate,
			&transaction.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan transaction history: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate transaction history: %w", err)
	}

	return transactions, nil
}
