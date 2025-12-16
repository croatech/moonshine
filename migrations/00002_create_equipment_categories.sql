-- +goose Up
-- +goose StatementBegin
CREATE TABLE equipment_categories (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS equipment_categories;
-- +goose StatementEnd

