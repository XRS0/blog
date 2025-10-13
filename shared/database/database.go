package database
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
