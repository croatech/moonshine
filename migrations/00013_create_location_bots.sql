-- +goose Up
-- +goose StatementBegin
CREATE TABLE location_bots (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    location_id INTEGER NOT NULL,
    bot_id INTEGER NOT NULL,
    CONSTRAINT fk_location_bots_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE,
    CONSTRAINT fk_location_bots_bot FOREIGN KEY (bot_id) REFERENCES bots(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS location_bots;
-- +goose StatementEnd

