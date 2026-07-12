package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Category struct {
	CategoryID int64
	Name       string
}

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) ResolveByName(ctx context.Context, name string) (*Category, error) {
	normalizedName := strings.ToLower(strings.TrimSpace(name))
	if normalizedName == "" {
		normalizedName = "lainnya"
	}

	query := `
		SELECT category_id, name
		FROM category
		WHERE account_id IS NULL
		  AND lower(name) = $1
		UNION ALL
		SELECT category_id, name
		FROM category
		WHERE account_id IS NULL
		  AND name = 'lainnya'
		  AND NOT EXISTS (
			SELECT 1 FROM category WHERE account_id IS NULL AND lower(name) = $1
		  )
		LIMIT 1
	`

	var category Category
	if err := r.db.QueryRowContext(ctx, query, normalizedName).Scan(&category.CategoryID, &category.Name); err != nil {
		return nil, fmt.Errorf("resolve category by name: %w", err)
	}

	return &category, nil
}
