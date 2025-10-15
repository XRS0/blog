package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

// NewDB creates a new Bun database connection
func NewDB(dsn string, logger *slog.Logger) (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	// Add query hook for debugging (only in development)
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	// Test connection
	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connected")
	return db, nil
}

// RunMigrations runs auto migrations for given models
func RunMigrations(ctx context.Context, db *bun.DB, models []interface{}, logger *slog.Logger) error {
	for _, model := range models {
		logger.Info("creating table", "model", fmt.Sprintf("%T", model))
		_, err := db.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create table for %T: %w", model, err)
		}
	}
	logger.Info("migrations completed", "tables", len(models))
	return nil
}

// AddColumn adds a column to a table if it doesn't exist
func AddColumn(ctx context.Context, db *bun.DB, table, column, columnType string, logger *slog.Logger) error {
	// Check if column exists
	var exists bool
	err := db.NewRaw(`
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = ? AND column_name = ?
		)
	`, table, column).Scan(ctx, &exists)
	
	if err != nil {
		return fmt.Errorf("failed to check column existence: %w", err)
	}
	
	if exists {
		logger.Info("column already exists", "table", table, "column", column)
		return nil
	}
	
	// Add column
	logger.Info("adding column", "table", table, "column", column, "type", columnType)
	_, err = db.ExecContext(ctx, fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, columnType))
	if err != nil {
		return fmt.Errorf("failed to add column: %w", err)
	}
	
	logger.Info("column added successfully", "table", table, "column", column)
	return nil
}
