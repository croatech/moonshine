-- +goose Up
-- +goose StatementBegin
CREATE TABLE tool_items (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL DEFAULT 0,
    required_skill INTEGER NOT NULL DEFAULT 0,
    tool_category_id INTEGER,
    image VARCHAR(255),
    CONSTRAINT fk_tool_items_category FOREIGN KEY (tool_category_id) REFERENCES tool_categories(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tool_items;
-- +goose StatementEnd

