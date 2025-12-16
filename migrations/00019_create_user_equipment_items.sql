-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_equipment_items (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id INTEGER NOT NULL,
    equipment_item_id INTEGER NOT NULL,
    CONSTRAINT fk_user_equipment_items_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_equipment_items_item FOREIGN KEY (equipment_item_id) REFERENCES equipment_items(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_equipment_items;
-- +goose StatementEnd

