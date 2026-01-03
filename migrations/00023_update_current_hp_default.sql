-- +goose Up
-- +goose StatementBegin

ALTER TABLE users ALTER COLUMN current_hp SET DEFAULT 20;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE users ALTER COLUMN current_hp SET DEFAULT 0;

-- +goose StatementEnd










