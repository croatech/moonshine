-- +goose Up
-- +goose StatementBegin
CREATE TABLE stuffs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    stuffable_type VARCHAR(255) NOT NULL,
    stuffable_id UUID NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_stuffs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stuffs;
-- +goose StatementEnd

