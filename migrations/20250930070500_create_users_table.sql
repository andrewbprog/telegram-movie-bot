-- +goose Up
CREATE TABLE IF NOT EXISTS users (
     id BIGSERIAL PRIMARY KEY,
     user_id TEXT NOT NULL UNIQUE,
     telegram_id TEXT NOT NULL UNIQUE,
     telegram_chat_id INT4 NOT NULL UNIQUE,
     created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
     updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
     );

-- +goose Down
DROP TABLE IF EXISTS users;
