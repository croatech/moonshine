-- +goose Up
-- +goose StatementBegin
CREATE TABLE fights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id UUID NOT NULL,
    bot_id UUID NOT NULL,
    status INTEGER NOT NULL DEFAULT 0,
    winner_type VARCHAR(255),
    dropped_gold INTEGER NOT NULL DEFAULT 0,
    winner_id UUID,
    dropped_item_id UUID,
    dropped_item_type VARCHAR(255),
    CONSTRAINT fk_fights_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_fights_bot FOREIGN KEY (bot_id) REFERENCES bots(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS fights;
-- +goose StatementEnd

