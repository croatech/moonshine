-- +goose Up
-- +goose StatementBegin
CREATE TYPE body_part AS ENUM ('HEAD', 'CHEST', 'BELT', 'LEGS', 'HANDS');

CREATE TABLE rounds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    fight_id UUID NOT NULL,
    player_damage INTEGER NOT NULL DEFAULT 0,
    bot_damage INTEGER NOT NULL DEFAULT 0,
    status fight_status NOT NULL DEFAULT 'IN_PROGRESS',
    player_hp INTEGER NOT NULL DEFAULT 0,
    bot_hp INTEGER NOT NULL DEFAULT 0,
    player_attack_point body_part,
    player_defense_point body_part,
    bot_attack_point body_part,
    bot_defense_point body_part,
    CONSTRAINT fk_rounds_fight FOREIGN KEY (fight_id) REFERENCES fights(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rounds;
DROP TYPE IF EXISTS body_part;
-- +goose StatementEnd
