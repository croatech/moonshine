-- +goose Up
-- +goose StatementBegin

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for all tables with updated_at column
CREATE TRIGGER update_avatars_updated_at BEFORE UPDATE ON avatars FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_equipment_categories_updated_at BEFORE UPDATE ON equipment_categories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_equipment_items_updated_at BEFORE UPDATE ON equipment_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tool_categories_updated_at BEFORE UPDATE ON tool_categories FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tool_items_updated_at BEFORE UPDATE ON tool_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_locations_updated_at BEFORE UPDATE ON locations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_bots_updated_at BEFORE UPDATE ON bots FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_resources_updated_at BEFORE UPDATE ON resources FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON events FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_messages_updated_at BEFORE UPDATE ON messages FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_movements_updated_at BEFORE UPDATE ON movements FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_location_bots_updated_at BEFORE UPDATE ON location_bots FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_location_locations_updated_at BEFORE UPDATE ON location_locations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_location_resources_updated_at BEFORE UPDATE ON location_resources FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_fights_updated_at BEFORE UPDATE ON fights FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_rounds_updated_at BEFORE UPDATE ON rounds FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_stuffs_updated_at BEFORE UPDATE ON stuffs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_equipment_items_updated_at BEFORE UPDATE ON user_equipment_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_tool_items_updated_at BEFORE UPDATE ON user_tool_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_bot_equipment_items_updated_at BEFORE UPDATE ON bot_equipment_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_avatars_updated_at ON avatars;
DROP TRIGGER IF EXISTS update_equipment_categories_updated_at ON equipment_categories;
DROP TRIGGER IF EXISTS update_equipment_items_updated_at ON equipment_items;
DROP TRIGGER IF EXISTS update_tool_categories_updated_at ON tool_categories;
DROP TRIGGER IF EXISTS update_tool_items_updated_at ON tool_items;
DROP TRIGGER IF EXISTS update_locations_updated_at ON locations;
DROP TRIGGER IF EXISTS update_bots_updated_at ON bots;
DROP TRIGGER IF EXISTS update_resources_updated_at ON resources;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_events_updated_at ON events;
DROP TRIGGER IF EXISTS update_messages_updated_at ON messages;
DROP TRIGGER IF EXISTS update_movements_updated_at ON movements;
DROP TRIGGER IF EXISTS update_location_bots_updated_at ON location_bots;
DROP TRIGGER IF EXISTS update_location_locations_updated_at ON location_locations;
DROP TRIGGER IF EXISTS update_location_resources_updated_at ON location_resources;
DROP TRIGGER IF EXISTS update_fights_updated_at ON fights;
DROP TRIGGER IF EXISTS update_rounds_updated_at ON rounds;
DROP TRIGGER IF EXISTS update_stuffs_updated_at ON stuffs;
DROP TRIGGER IF EXISTS update_user_equipment_items_updated_at ON user_equipment_items;
DROP TRIGGER IF EXISTS update_user_tool_items_updated_at ON user_tool_items;
DROP TRIGGER IF EXISTS update_bot_equipment_items_updated_at ON bot_equipment_items;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- +goose StatementEnd




