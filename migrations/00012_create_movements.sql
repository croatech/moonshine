-- +goose Up
-- +goose StatementBegin
CREATE TYPE movement_status AS ENUM ('active', 'finished');

CREATE TABLE movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id UUID NOT NULL,
    status movement_status NOT NULL DEFAULT 'active',
    CONSTRAINT fk_movements_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_movements_user_id ON movements(user_id);
CREATE INDEX idx_movements_status ON movements(status);
CREATE INDEX idx_movements_user_status ON movements(user_id, status) WHERE deleted_at IS NULL;

CREATE TABLE movements_cells (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movement_id UUID NOT NULL,
    from_cell_id UUID NOT NULL,
    to_cell_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_movements_cells_movement FOREIGN KEY (movement_id) REFERENCES movements(id) ON DELETE CASCADE,
    CONSTRAINT fk_movements_cells_from FOREIGN KEY (from_cell_id) REFERENCES locations(id) ON DELETE CASCADE,
    CONSTRAINT fk_movements_cells_to FOREIGN KEY (to_cell_id) REFERENCES locations(id) ON DELETE CASCADE
);

CREATE INDEX idx_movements_cells_movement_id ON movements_cells(movement_id);
CREATE INDEX idx_movements_cells_created_at ON movements_cells(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS movements_cells;
DROP TABLE IF EXISTS movements;
DROP TYPE IF EXISTS movement_status;
-- +goose StatementEnd

