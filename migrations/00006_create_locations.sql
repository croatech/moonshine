-- +goose Up
-- +goose StatementBegin
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255),
    cell BOOLEAN NOT NULL DEFAULT true,
    inactive BOOLEAN NOT NULL DEFAULT true,
    parent_id INTEGER,
    CONSTRAINT fk_locations_parent FOREIGN KEY (parent_id) REFERENCES locations(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS locations;
-- +goose StatementEnd

