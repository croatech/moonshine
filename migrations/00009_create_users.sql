-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    avatar_id UUID,
    location_id UUID NOT NULL,
    armor_slot INTEGER NOT NULL DEFAULT 0,
    attack INTEGER NOT NULL DEFAULT 1,
    defense INTEGER NOT NULL DEFAULT 1,
    belt_slot INTEGER NOT NULL DEFAULT 0,
    bracers_slot INTEGER NOT NULL DEFAULT 0,
    cloak_slot INTEGER NOT NULL DEFAULT 0,
    current_hp INTEGER NOT NULL DEFAULT 0,
    exp INTEGER NOT NULL DEFAULT 0,
    exp_next INTEGER NOT NULL DEFAULT 100,
    fishing_skill INTEGER NOT NULL DEFAULT 0,
    fishing_slot INTEGER NOT NULL DEFAULT 0,
    foots_slot INTEGER NOT NULL DEFAULT 0,
    free_stats INTEGER NOT NULL DEFAULT 10,
    gloves_slot INTEGER NOT NULL DEFAULT 0,
    gold INTEGER NOT NULL DEFAULT 0,
    helmet_slot INTEGER NOT NULL DEFAULT 0,
    hp INTEGER NOT NULL DEFAULT 20,
    level INTEGER NOT NULL DEFAULT 1,
    lumberjacking_skill INTEGER NOT NULL DEFAULT 0,
    lumberjacking_slot INTEGER NOT NULL DEFAULT 0,
    mail_slot INTEGER NOT NULL DEFAULT 0,
    necklace_slot INTEGER NOT NULL DEFAULT 0,
    pants_slot INTEGER NOT NULL DEFAULT 0,
    ring_slot INTEGER NOT NULL DEFAULT 0,
    shield_slot INTEGER NOT NULL DEFAULT 0,
    weapon_slot INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT fk_users_avatar FOREIGN KEY (avatar_id) REFERENCES avatars(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE RESTRICT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

