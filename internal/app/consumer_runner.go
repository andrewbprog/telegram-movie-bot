package app

import (
	"context"
	"go.uber.org/zap"
	"tlgbs/internal/infrastructure/kafka"
	"tlgbs/internal/service"
)

func RunKafkaConsumer(ctx context.Context, brokers []string, topic, groupID string, notifSvc *service.NotificationService, logger *zap.Logger) {
	go func() {
		logger.Info("starting kafka consumer...",
			zap.String("topic", topic),
			zap.String("group_id", groupID),
			zap.Strings("brokers", brokers),
		)

		if err := kafka.StartConsumer(ctx, brokers, topic, groupID, notifSvc, logger); err != nil {
			logger.Error("kafka consumer failed", zap.Error(err))
		}
	}()
}
