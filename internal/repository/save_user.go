package repository

import (
	"context"
	"fmt"
	"telegram-movie-bot/internal/repository/models"
)

const (
	queryInsert = `
		INSERT INTO users (user_id, telegram_id, telegram_chat_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (telegram_id) 
		DO UPDATE SET telegram_chat_id = EXCLUDED.telegram_chat_id;
	`
)

func (r *UserRepository) SaveUser(ctx context.Context, user *models.User) error {

	_, err := r.db.Exec(ctx, queryInsert, user.UserID, user.TelegramID, user.TelegramChatID)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}
