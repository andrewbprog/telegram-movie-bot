package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

func RunMigrations(ctx context.Context, migrationsDir, dbURL string, logger *zap.Logger) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}
	defer func(db *sql.DB) {
		err2 := db.Close()
		if err2 != nil {
			logger.Error("failed to close database", zap.Error(err2))
		}
	}(db)

	goose.SetBaseFS(nil) // если миграции в локальной папке, это не нужно менять

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		err := goose.Up(db, migrationsDir)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("migrations timeout: %w", ctx.Err())
	case err := <-done:
		if err != nil {
			logger.Error("goose migrations failed", zap.Error(err))
			return fmt.Errorf("migrations failed: %w", err)
		}
		logger.Info("migrations applied successfully with goose")
		return nil
	}
}
