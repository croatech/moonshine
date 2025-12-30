-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    avatar_id UUID,
    location_id UUID NOT NULL,
    attack INTEGER NOT NULL DEFAULT 1,
    defense INTEGER NOT NULL DEFAULT 1,
    current_hp INTEGER NOT NULL DEFAULT 20,
    exp INTEGER NOT NULL DEFAULT 0,
    fishing_skill INTEGER NOT NULL DEFAULT 0,
    fishing_slot INTEGER NOT NULL DEFAULT 0,
    free_stats INTEGER NOT NULL DEFAULT 15,
    gold INTEGER NOT NULL DEFAULT 0,
    hp INTEGER NOT NULL DEFAULT 20,
    level INTEGER NOT NULL DEFAULT 1,
    lumberjacking_skill INTEGER NOT NULL DEFAULT 0,
    lumberjacking_slot INTEGER NOT NULL DEFAULT 0,
    chest_equipment_item_id UUID,
    belt_equipment_item_id UUID,
    head_equipment_item_id UUID,
    neck_equipment_item_id UUID,
    weapon_equipment_item_id UUID,
    shield_equipment_item_id UUID,
    legs_equipment_item_id UUID,
    feet_equipment_item_id UUID,
    arms_equipment_item_id UUID,
    hands_equipment_item_id UUID,
    ring1_equipment_item_id UUID,
    ring2_equipment_item_id UUID,
    ring3_equipment_item_id UUID,
    ring4_equipment_item_id UUID,
    CONSTRAINT fk_users_avatar FOREIGN KEY (avatar_id) REFERENCES avatars(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE RESTRICT,
    CONSTRAINT fk_users_chest_equipment FOREIGN KEY (chest_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_belt_equipment FOREIGN KEY (belt_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_head_equipment FOREIGN KEY (head_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_neck_equipment FOREIGN KEY (neck_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_weapon_equipment FOREIGN KEY (weapon_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_shield_equipment FOREIGN KEY (shield_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_legs_equipment FOREIGN KEY (legs_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_feet_equipment FOREIGN KEY (feet_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_arms_equipment FOREIGN KEY (arms_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_hands_equipment FOREIGN KEY (hands_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_ring1_equipment FOREIGN KEY (ring1_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_ring2_equipment FOREIGN KEY (ring2_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_ring3_equipment FOREIGN KEY (ring3_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL,
    CONSTRAINT fk_users_ring4_equipment FOREIGN KEY (ring4_equipment_item_id) REFERENCES equipment_items(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

