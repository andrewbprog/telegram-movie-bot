package service

import (
	"context"
	"go.uber.org/zap"
	tgclient "tlgbs/internal/infrastructure/telegram-client"
	"tlgbs/internal/repository"
)

const (
	msg = "📢 Ваша лента рекомендаций обновилась!\nВведите /recommendations, чтобы посмотреть новые рекомендации."
)

type NotificationService struct {
	repo     *repository.UserRepository
	telegram *tgclient.TgClient
	logger   *zap.Logger
}

func NewNotificationService(repo *repository.UserRepository, tgClient *tgclient.TgClient, logger *zap.Logger) *NotificationService {
	return &NotificationService{
		repo:     repo,
		telegram: tgClient,
		logger:   logger,
	}
}

func (s *NotificationService) HandleRecs(ctx context.Context, userID string) error {
	chatID, err := s.repo.GetChatIDByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("chat not found for user", zap.String("user_id", userID), zap.Error(err))
		return nil
	}

	if err := s.telegram.SendMessage(chatID, msg); err != nil {
		s.logger.Error("failed to send telegram message",
			zap.String("user_id", userID),
			zap.Int64("chat_id", chatID),
			zap.Error(err),
		)
		return nil
	}

	s.logger.Info("notification sent", zap.String("user_id", userID), zap.Int64("chat_id", chatID))
	return nil
}
