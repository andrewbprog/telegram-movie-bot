package migrate

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func RunMigrations(ctx context.Context, migrationsPath, dbURL string, logger *zap.Logger) error {
	migURL := fmt.Sprintf("file://%s", migrationsPath)
	m, err := migrate.New(migURL, dbURL)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		err := m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			done <- err
			return
		}
		done <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		if err != nil {
			logger.Error("migrate up failed", zap.Error(err))
			return fmt.Errorf("error: %w", err)
		}
		logger.Info("migrations applied successfully")
		return nil
	}
}
