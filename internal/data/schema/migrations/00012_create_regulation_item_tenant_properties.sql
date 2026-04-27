-- +goose Up
-- +goose StatementBegin
CREATE TABLE regulation_item_properties (
    regulation_item_id UUID NOT NULL REFERENCES regulation_items(id) ON DELETE CASCADE,
    property_id UUID NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
    PRIMARY KEY (regulation_item_id, property_id)
);

-- Migrate existing data
INSERT INTO regulation_item_properties (regulation_item_id, property_id)
SELECT id, tenant_property_id FROM regulation_items WHERE tenant_property_id IS NOT NULL;

-- Remove old column
ALTER TABLE regulation_items DROP COLUMN tenant_property_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE regulation_items ADD COLUMN tenant_property_id UUID REFERENCES tenants_properties(id) ON DELETE SET NULL;

-- Note: This back-migration only restores one property per item (the first one found)
UPDATE regulation_items ri
SET tenant_property_id = tp.id
FROM regulation_item_properties rip
JOIN tenants_properties tp ON rip.property_id = tp.property_id
WHERE ri.id = rip.regulation_item_id;

DROP TABLE regulation_item_properties;
-- +goose StatementEnd
