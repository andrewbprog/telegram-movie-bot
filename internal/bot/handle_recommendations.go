package bot

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

func (b *Bot) handleRecommendations(ctx context.Context, chatID int64, userID string) error {
	recs, err := b.gw.GetRecommendations(ctx, userID)
	if err != nil {
		b.logger.Error("failed to get recommendations", zap.Error(err))

		if sendErr := b.SendAnswerToTelegram(chatID, "❌ Не удалось получить рекомендации, попробуйте позже."); sendErr != nil {
			b.logger.Error("failed to get recommendations", zap.Error(sendErr))
		}
		return nil
	}

	if len(recs) == 0 {
		if sendErr := b.SendAnswerToTelegram(chatID, "Пока нет рекомендаций для вас."); sendErr != nil {
			b.logger.Error("failed to send empty recommendations message", zap.Error(sendErr))
		}
		return nil
	}

	for _, r := range recs {
		msgText := fmt.Sprintf("🎬 %s\n%s", r.Title, r.Description)
		if sendErr := b.SendAnswerToTelegram(chatID, msgText); sendErr != nil {
			b.logger.Error("failed to send recommendation", zap.Error(sendErr))
		}
	}
	return nil
}
