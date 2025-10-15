package repository

import (
	"context"
	"fmt"
	"tlgbs/internal/repository/models"
)

const (
	InsertQuery = `
		INSERT INTO users (user_id, telegram_id, telegram_chat_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (telegram_id) 
		DO UPDATE SET telegram_chat_id = EXCLUDED.telegram_chat_id;
	`
)

func (r *UserRepository) SaveUser(ctx context.Context, user *models.User) error {

	_, err := r.db.Exec(ctx, InsertQuery, user.UserID, user.TelegramID, user.TelegramChatID)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}
