package models

type User struct {
	ID             int64  `db:"id"`
	UserID         string `db:"user_id"`
	TelegramID     string `db:"telegram_id"`
	TelegramChatID int64  `db:"telegram_chat_id"`
}
