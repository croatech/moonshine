-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id INTEGER NOT NULL,
    recipient_id INTEGER NOT NULL,
    text TEXT,
    CONSTRAINT fk_messages_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd

