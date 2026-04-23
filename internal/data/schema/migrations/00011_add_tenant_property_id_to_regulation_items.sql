-- +goose Up
-- +goose StatementBegin
ALTER TABLE regulation_items ADD COLUMN tenant_property_id UUID REFERENCES tenants_properties(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE regulation_items DROP COLUMN tenant_property_id;
-- +goose StatementEnd
