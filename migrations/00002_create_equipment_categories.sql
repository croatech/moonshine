-- +goose Up
-- +goose StatementBegin
CREATE TYPE equipment_category_type AS ENUM (
    'chest',
    'belt',
    'head',
    'neck',
    'weapon',
    'shield',
    'legs',
    'feet',
    'arms',
    'hands',
    'ring'
);

CREATE TABLE equipment_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    type equipment_category_type
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS equipment_categories;
DROP TYPE IF EXISTS equipment_category_type;
-- +goose StatementEnd

