-- +goose Up
-- +goose StatementBegin
ALTER TABLE equipment_items ADD COLUMN slug VARCHAR(255);
CREATE INDEX idx_equipment_items_slug ON equipment_items(slug) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_equipment_items_slug;
ALTER TABLE equipment_items DROP COLUMN IF EXISTS slug;
-- +goose StatementEnd

