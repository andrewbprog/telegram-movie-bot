package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"tlgbs/internal/infrastructure/models"
	"tlgbs/internal/service"
)

func StartConsumer(ctx context.Context, brokers []string, topic, groupID string, notifySvc *service.NotificationService, logger *zap.Logger) error {

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})
	defer r.Close()

	logger.Info("kafka consumer started", zap.String("topic", topic))

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				logger.Info("consumer context canceled, stopping gracefully")
				return nil
			}
			logger.Error("failed to read message", zap.Error(err))
			continue
		}

		var evt models.RecommendationEvent
		if err := json.Unmarshal(m.Value, &evt); err != nil {
			logger.Error("invalid kafka message format", zap.Error(err))
			continue
		}

		if evt.UserID == "" {
			logger.Error("received empty user_id, skipping message")
			continue
		}

		if err := notifySvc.HandleRecs(ctx, evt.UserID); err != nil {
			logger.Error("failed to handle notification", zap.Error(err))
		}

		logger.Info("received user", zap.String("user_id", evt.UserID))
	}
}
