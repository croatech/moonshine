-- +goose Up
-- +goose StatementBegin
ALTER TABLE fights ADD COLUMN exp INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE fights DROP COLUMN IF EXISTS exp;
-- +goose StatementEnd
