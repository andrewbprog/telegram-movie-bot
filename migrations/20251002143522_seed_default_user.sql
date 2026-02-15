-- +goose Up
INSERT INTO users (user_id, telegram_id, telegram_chat_id, created_at, updated_at)
VALUES ('freeman04221', 'freeman04221', 84043246, now(), now())
    ON CONFLICT (user_id) DO NOTHING;

-- +goose Down
DELETE FROM users WHERE user_id = 'freeman04221';

