package repository

import (
	"context"
	"fmt"
)

const (
	queryGetUserId = `SELECT user_id FROM users WHERE telegram_id = $1`
)

func (r *UserRepository) GetUserIDByTelegramID(ctx context.Context, telegramID string) (string, error) {

	var userID string
	err := r.db.QueryRow(ctx, queryGetUserId, telegramID).Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user_id by telegram_id: %w", err)
	}

	return userID, nil
}
