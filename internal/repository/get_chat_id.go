package repository

import (
	"context"
	"fmt"
)

const (
	queryGetChatId = `SELECT telegram_chat_id FROM users WHERE user_id = $1`
)

func (r *UserRepository) GetChatIDByUserID(ctx context.Context, userID string) (int64, error) {

	var chatID int64
	err := r.db.QueryRow(ctx, queryGetChatId, userID).Scan(&chatID)
	if err != nil {
		return 0, fmt.Errorf("failed to get telegram_chat_id by user_id: %w", err)
	}

	return chatID, nil
}
