-- +goose Up
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS fishing_skill;
ALTER TABLE users DROP COLUMN IF EXISTS fishing_slot;
ALTER TABLE users DROP COLUMN IF EXISTS lumberjacking_skill;
ALTER TABLE users DROP COLUMN IF EXISTS lumberjacking_slot;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN fishing_skill INTEGER NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN fishing_slot INTEGER NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN lumberjacking_skill INTEGER NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN lumberjacking_slot INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd
