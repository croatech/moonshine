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

-- Function to update user stats when equipment changes
CREATE OR REPLACE FUNCTION update_user_stats_from_equipment()
RETURNS TRIGGER AS $$
DECLARE
    total_attack INTEGER := 1;
    total_defense INTEGER := 1;
    total_hp INTEGER := 20;
    equipment_changed BOOLEAN := FALSE;
BEGIN
    -- Check if this is INSERT or if equipment fields changed
    IF TG_OP = 'INSERT' THEN
        equipment_changed := TRUE;
    ELSE
        -- Check if any equipment field changed
        equipment_changed := (
            NEW.chest_equipment_item_id IS DISTINCT FROM OLD.chest_equipment_item_id OR
            NEW.belt_equipment_item_id IS DISTINCT FROM OLD.belt_equipment_item_id OR
            NEW.head_equipment_item_id IS DISTINCT FROM OLD.head_equipment_item_id OR
            NEW.neck_equipment_item_id IS DISTINCT FROM OLD.neck_equipment_item_id OR
            NEW.weapon_equipment_item_id IS DISTINCT FROM OLD.weapon_equipment_item_id OR
            NEW.shield_equipment_item_id IS DISTINCT FROM OLD.shield_equipment_item_id OR
            NEW.legs_equipment_item_id IS DISTINCT FROM OLD.legs_equipment_item_id OR
            NEW.feet_equipment_item_id IS DISTINCT FROM OLD.feet_equipment_item_id OR
            NEW.arms_equipment_item_id IS DISTINCT FROM OLD.arms_equipment_item_id OR
            NEW.hands_equipment_item_id IS DISTINCT FROM OLD.hands_equipment_item_id OR
            NEW.ring1_equipment_item_id IS DISTINCT FROM OLD.ring1_equipment_item_id OR
            NEW.ring2_equipment_item_id IS DISTINCT FROM OLD.ring2_equipment_item_id OR
            NEW.ring3_equipment_item_id IS DISTINCT FROM OLD.ring3_equipment_item_id OR
            NEW.ring4_equipment_item_id IS DISTINCT FROM OLD.ring4_equipment_item_id
        );
    END IF;

    -- Only recalculate if equipment changed
    IF equipment_changed THEN
        -- Calculate base stats
        total_attack := 1;
        total_defense := 1;
        total_hp := 20;

        -- Sum stats from all equipment items (only non-null IDs)
        SELECT 
            1 + COALESCE(SUM(ei.attack), 0),
            1 + COALESCE(SUM(ei.defense), 0),
            20 + COALESCE(SUM(ei.hp), 0)
        INTO total_attack, total_defense, total_hp
        FROM equipment_items ei
        WHERE ei.deleted_at IS NULL
            AND ei.id IN (
                NEW.chest_equipment_item_id,
                NEW.belt_equipment_item_id,
                NEW.head_equipment_item_id,
                NEW.neck_equipment_item_id,
                NEW.weapon_equipment_item_id,
                NEW.shield_equipment_item_id,
                NEW.legs_equipment_item_id,
                NEW.feet_equipment_item_id,
                NEW.arms_equipment_item_id,
                NEW.hands_equipment_item_id,
                NEW.ring1_equipment_item_id,
                NEW.ring2_equipment_item_id,
                NEW.ring3_equipment_item_id,
                NEW.ring4_equipment_item_id
            );

        -- Update user stats
        NEW.attack := total_attack;
        NEW.defense := total_defense;
        NEW.hp := total_hp;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger before insert or update on users
-- Optimized: only recalculates stats when equipment actually changes
CREATE TRIGGER trigger_update_user_stats_from_equipment
    BEFORE INSERT OR UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_user_stats_from_equipment();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_update_user_stats_from_equipment ON users;
DROP FUNCTION IF EXISTS update_user_stats_from_equipment();
DROP TABLE IF EXISTS users;
-- +goose StatementEnd

