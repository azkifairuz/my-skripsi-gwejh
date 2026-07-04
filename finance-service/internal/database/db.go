package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed schema.sql seed.sql
var migrationFiles embed.FS

func Open(ctx context.Context, config Config) (*sql.DB, error) {
	if err := ensureDatabase(ctx, config); err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

func ensureDatabase(ctx context.Context, config Config) error {
	db, err := sql.Open("pgx", config.MaintenanceDSN())
	if err != nil {
		return fmt.Errorf("open maintenance database: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping maintenance database: %w", err)
	}

	var exists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)", config.Name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check database exists: %w", err)
	}

	if exists {
		return nil
	}

	if _, err := db.ExecContext(ctx, `CREATE DATABASE `+quoteIdentifier(config.Name)); err != nil {
		return fmt.Errorf("create database %s: %w", config.Name, err)
	}

	return nil
}

func quoteIdentifier(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}

func Migrate(ctx context.Context, db *sql.DB) error {
	if err := execEmbeddedSQL(ctx, db, "schema.sql"); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}

	if err := execEmbeddedSQL(ctx, db, "seed.sql"); err != nil {
		return fmt.Errorf("apply seed: %w", err)
	}

	return nil
}

func execEmbeddedSQL(ctx context.Context, db *sql.DB, name string) error {
	query, err := migrationFiles.ReadFile(name)
	if err != nil {
		return fmt.Errorf("read %s: %w", name, err)
	}

	if _, err := db.ExecContext(ctx, string(query)); err != nil {
		return fmt.Errorf("execute %s: %w", name, err)
	}

	return nil
}
