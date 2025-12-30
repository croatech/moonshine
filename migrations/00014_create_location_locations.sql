-- +goose Up
-- +goose StatementBegin
CREATE TABLE location_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    location_id UUID NOT NULL,
    near_location_id UUID NOT NULL,
    CONSTRAINT fk_location_locations_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE,
    CONSTRAINT fk_location_locations_near FOREIGN KEY (near_location_id) REFERENCES locations(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS location_locations;
-- +goose StatementEnd

