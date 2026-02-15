package bot

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

func (b *Bot) handleRandomMovie(ctx context.Context, chatID int64, userID string) error {

	movie, err := b.gw.GetRandomMovie(ctx, userID)
	if err != nil {
		b.logger.Error("failed to get random movie", zap.Error(err))

		if sendErr := b.SendAnswerToTelegram(chatID, "❌ Не удалось получить случайный фильм. Попробуйте позже."); sendErr != nil {
			b.logger.Error("failed to send error message to user", zap.Error(sendErr))
		}
		return nil
	}

	if movie == nil {
		if sendErr := b.SendAnswerToTelegram(chatID, "😕 Фильм не найден, попробуйте позже."); sendErr != nil {
			b.logger.Error("failed to send empty movie message", zap.Error(sendErr))
		}
		return nil
	}

	text := fmt.Sprintf("🎬 %s\n\n%s", movie.Title, movie.Description)

	if sendErr := b.SendAnswerToTelegram(chatID, text); sendErr != nil {
		b.logger.Error("failed to send random movie", zap.Error(sendErr))
		return nil
	}

	b.logger.Info("random movie sent",
		zap.String("user_id", userID),
		zap.Int64("chat_id", chatID),
		zap.String("title", movie.Title),
	)
	return nil
}
