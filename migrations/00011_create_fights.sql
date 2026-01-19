-- +goose Up
-- +goose StatementBegin
CREATE TYPE fight_status AS ENUM ('IN_PROGRESS', 'FINISHED');

CREATE TABLE fights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id UUID NOT NULL,
    bot_id UUID NOT NULL,
    status fight_status NOT NULL DEFAULT 'IN_PROGRESS',
    user_won BOOLEAN NOT NULL DEFAULT false,
    dropped_gold INTEGER NOT NULL DEFAULT 0,
    dropped_item_id UUID,
    CONSTRAINT fk_fights_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_fights_bot FOREIGN KEY (bot_id) REFERENCES bots(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS fights;
DROP TYPE IF EXISTS fight_status;
-- +goose StatementEnd
