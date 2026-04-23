-- +goose Up
-- +goose StatementBegin
CREATE TABLE regulation_item_tenant_properties (
    regulation_item_model_id UUID NOT NULL REFERENCES regulation_items(id) ON DELETE CASCADE,
    tenant_property_model_id UUID NOT NULL REFERENCES tenants_properties(id) ON DELETE CASCADE,
    PRIMARY KEY (regulation_item_model_id, tenant_property_model_id)
);

-- Migrate existing data
INSERT INTO regulation_item_tenant_properties (regulation_item_model_id, tenant_property_model_id)
SELECT id, tenant_property_id FROM regulation_items WHERE tenant_property_id IS NOT NULL;

-- Remove old column
ALTER TABLE regulation_items DROP COLUMN tenant_property_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE regulation_items ADD COLUMN tenant_property_id UUID REFERENCES tenants_properties(id) ON DELETE SET NULL;

-- Note: This back-migration only restores one property per item (the first one found)
UPDATE regulation_items ri
SET tenant_property_id = ritp.tenant_property_model_id
FROM regulation_item_tenant_properties ritp
WHERE ri.id = ritp.regulation_item_model_id;

DROP TABLE regulation_item_tenant_properties;
-- +goose StatementEnd
