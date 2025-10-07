package app

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"tlgbs/config"
	"tlgbs/internal/bot"
	"tlgbs/internal/migrate"
	"tlgbs/internal/postgres"
)

func Start(ctx context.Context) error {

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, err := zap.NewProduction()
	if err != nil {
		logger.Error("failed to init zap logger", zap.Error(err))
		return fmt.Errorf("logger error: %w", err)
	}

	cfg, err := config.New()
	if err != nil {
		logger.Error("init configuration failed", zap.Error(err))
		return fmt.Errorf("config error: %w", err)
	}

	migCtx, migCancel := context.WithTimeout(ctx, cfg.MigrationTimeout)
	defer migCancel()

	if err := migrate.RunMigrations(migCtx, "migrations", cfg.Database.DatabaseURL(), logger); err != nil {
		logger.Error("migrations failed", zap.Error(err))
		return fmt.Errorf("migrations error: %w", err)
	}

	pool, err := postgres.NewPool(ctx, cfg.Database.DatabaseURL())
	if err != nil {
		logger.Error("db connection failed", zap.Error(err))
		return fmt.Errorf("pool error: %w", err)
	}
	defer pool.Close()

	b, err := bot.NewBot(cfg.TelegramToken, logger)
	if err != nil {
		logger.Error("failed to init telegram bot", zap.Error(err))
		return fmt.Errorf("bot error: %w", err)
	}
	go b.Run(ctx)
	logger.Info("telegram bot started")

	<-ctx.Done()
	logger.Info("context cancelled, shutting down app")

	return nil
}
