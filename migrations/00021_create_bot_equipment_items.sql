-- +goose Up
-- +goose StatementBegin
CREATE TABLE bot_equipment_items (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    bot_id INTEGER NOT NULL,
    equipment_item_id INTEGER NOT NULL,
    CONSTRAINT fk_bot_equipment_items_bot FOREIGN KEY (bot_id) REFERENCES bots(id) ON DELETE CASCADE,
    CONSTRAINT fk_bot_equipment_items_item FOREIGN KEY (equipment_item_id) REFERENCES equipment_items(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bot_equipment_items;
-- +goose StatementEnd

