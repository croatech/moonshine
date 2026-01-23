-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'NECK' AND enumtypid = (SELECT oid FROM pg_type WHERE typname = 'body_part')) THEN
        ALTER TYPE body_part ADD VALUE 'NECK';
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Cannot remove enum value in PostgreSQL, so we recreate the type
CREATE TYPE body_part_new AS ENUM ('HEAD', 'CHEST', 'BELT', 'LEGS', 'HANDS');

ALTER TABLE rounds 
  ALTER COLUMN player_attack_point TYPE body_part_new USING player_attack_point::text::body_part_new,
  ALTER COLUMN player_defense_point TYPE body_part_new USING player_defense_point::text::body_part_new,
  ALTER COLUMN bot_attack_point TYPE body_part_new USING bot_attack_point::text::body_part_new,
  ALTER COLUMN bot_defense_point TYPE body_part_new USING bot_defense_point::text::body_part_new;

DROP TYPE body_part;
ALTER TYPE body_part_new RENAME TO body_part;
-- +goose StatementEnd
