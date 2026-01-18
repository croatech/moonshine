-- +goose Up
-- +goose StatementBegin
ALTER TABLE bots ADD COLUMN slug VARCHAR(255);
CREATE INDEX idx_bots_slug ON bots(slug) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_bots_slug;
ALTER TABLE bots DROP COLUMN IF EXISTS slug;
-- +goose StatementEnd
