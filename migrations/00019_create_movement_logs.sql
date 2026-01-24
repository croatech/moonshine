-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS movement_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    from_cell VARCHAR(255),
    to_cell VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_movement_logs_user_id ON movement_logs(user_id);
CREATE INDEX idx_movement_logs_created_at ON movement_logs(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS movement_logs;
-- +goose StatementEnd
