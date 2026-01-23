-- +goose Up
-- +goose StatementBegin
ALTER TABLE fights DROP COLUMN IF EXISTS user_won;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE fights ADD COLUMN user_won BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd
