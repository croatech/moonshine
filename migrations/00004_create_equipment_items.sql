-- +goose Up
-- +goose StatementBegin
CREATE TABLE equipment_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255),
    attack INTEGER NOT NULL DEFAULT 0,
    defense INTEGER NOT NULL DEFAULT 0,
    hp INTEGER NOT NULL DEFAULT 0,
    required_level INTEGER NOT NULL DEFAULT 1,
    price INTEGER NOT NULL DEFAULT 0,
    artifact BOOLEAN NOT NULL DEFAULT false,
    equipment_category_id UUID,
    image VARCHAR(255),
    CONSTRAINT fk_equipment_items_category FOREIGN KEY (equipment_category_id) REFERENCES equipment_categories(id) ON DELETE SET NULL
);

CREATE INDEX idx_equipment_items_slug ON equipment_items(slug) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_equipment_items_slug_unique ON equipment_items(slug) WHERE deleted_at IS NULL AND slug IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS equipment_items;
-- +goose StatementEnd

