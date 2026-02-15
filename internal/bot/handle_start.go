package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"tlgbs/internal/repository/models"
)

const msgText = "👋 Привет! Я бот для рекомендации фильмов!\n\n" +
	"Вот что я умею:\n" +
	"🎲 /random_movie — получить случайный фильм\n" +
	"🎬 /recommendations — показать рекомендации\n"

func (b *Bot) handleStart(chatID int64) error {
	ctx := context.Background()
	userID := uuid.New().String()

	newUser := &models.User{
		UserID:         userID,
		TelegramID:     userID,
		TelegramChatID: chatID,
	}

	if err := b.repo.SaveUser(ctx, newUser); err != nil {
		if sendErr := b.SendAnswerToTelegram(chatID, "❌ Произошла ошибка при добавлении пользователя. Попробуйте позже."); sendErr != nil {
			b.logger.Error("failed to save user", zap.Error(sendErr))
		}
		return fmt.Errorf("failed to save user %w", err)
	}

	msg := tgbotapi.NewMessage(chatID, msgText)
	if _, err := b.api.Send(msg); err != nil {
		return fmt.Errorf("failed to send start message %w", err)
	}

	b.logger.Info("user successfully started bot", zap.String("user_id", userID), zap.Int64("chat_id", chatID))
	return nil
}
