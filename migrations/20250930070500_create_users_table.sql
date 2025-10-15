-- +goose Up
CREATE TABLE IF NOT EXISTS users (
     id SERIAL PRIMARY KEY,
     user_id TEXT NOT NULL,
     telegram_id INTEGER UNIQUE NOT NULL,
     telegram_chat_id INTEGER NOT NULL,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
     );

-- +goose Down
DROP TABLE IF EXISTS users;