-- +goose Up
-- +goose StatementBegin
CREATE TABLE rounds (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    fight_id INTEGER NOT NULL,
    player_damage INTEGER NOT NULL DEFAULT 0,
    bot_damage INTEGER NOT NULL DEFAULT 0,
    status INTEGER NOT NULL DEFAULT 0,
    player_hp INTEGER NOT NULL DEFAULT 0,
    bot_hp INTEGER NOT NULL DEFAULT 0,
    player_attack_point VARCHAR(255),
    player_defense_point VARCHAR(255),
    bot_attack_point VARCHAR(255),
    bot_defense_point VARCHAR(255),
    CONSTRAINT fk_rounds_fight FOREIGN KEY (fight_id) REFERENCES fights(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rounds;
-- +goose StatementEnd

