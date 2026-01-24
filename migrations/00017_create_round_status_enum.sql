-- +goose Up
-- +goose StatementBegin
CREATE TYPE round_status AS ENUM ('IN_PROGRESS', 'FINISHED');

ALTER TABLE rounds 
  ALTER COLUMN status DROP DEFAULT;

ALTER TABLE rounds 
  ALTER COLUMN status TYPE round_status USING status::text::round_status;

ALTER TABLE rounds 
  ALTER COLUMN status SET DEFAULT 'IN_PROGRESS'::round_status;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE rounds 
  ALTER COLUMN status DROP DEFAULT;

ALTER TABLE rounds 
  ALTER COLUMN status TYPE fight_status USING status::text::fight_status;

ALTER TABLE rounds 
  ALTER COLUMN status SET DEFAULT 'IN_PROGRESS'::fight_status;

DROP TYPE IF EXISTS round_status;
-- +goose StatementEnd
