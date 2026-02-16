package app

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"telegram-movie-bot/config"
	"telegram-movie-bot/internal/bot"
	"telegram-movie-bot/internal/gateway"
	"telegram-movie-bot/internal/infrastructure/kafka"
	tgclient "telegram-movie-bot/internal/infrastructure/telegram-client"
	"telegram-movie-bot/internal/migrate"
	"telegram-movie-bot/internal/postgres"
	"telegram-movie-bot/internal/repository"
	"telegram-movie-bot/internal/service"
)

func Start(ctx context.Context) error {

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to init zap logger: %w", err)
	}

	cfg, err := config.New()
	if err != nil {
		logger.Error("failed to load configuration", zap.Error(err))
		return fmt.Errorf("failed to load config: %w", err)
	}

	migCtx, migCancel := context.WithTimeout(ctx, cfg.MigrationTimeout)
	defer migCancel()

	if err := migrate.RunMigrations(migCtx, "migrations", cfg.Database.DatabaseURL(), logger); err != nil {
		logger.Error("database migrations failed", zap.Error(err))
		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	pool, err := postgres.NewPool(ctx, cfg.Database.DatabaseURL())
	if err != nil {
		logger.Error("failed to connect to PostgreSQL", zap.Error(err))
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer pool.Close()
	repo := repository.NewUserRepository(pool)
	tg := tgclient.NewTgClient(cfg.TelegramToken)
	notifSvc := service.NewNotificationService(repo, tg, logger)

	// --- Запуск Consumer ---
	go func() {
		if err := kafka.StartConsumer(ctx, cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID, notifSvc, logger); err != nil {
			logger.Error("kafka consumer failed", zap.Error(err))
		}
	}()

	// --- Инициализация Telegram-бота ---
	gwClient := gateway.NewClient(cfg.GatewayURL, cfg.GatewayTimeout)
	b, err := bot.NewBot(cfg.TelegramToken, repo, gwClient, logger)
	if err != nil {
		logger.Error("failed to initialize telegram bot", zap.Error(err))
		return fmt.Errorf("failed to initialize telegram bot: %w", err)
	}

	// Запуск Telegram-бота
	go b.BotRun(ctx)

	<-ctx.Done()
	logger.Info("context cancelled, shutting down app")
	return nil
}
