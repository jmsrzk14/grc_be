-- +goose Up
-- +goose StatementBegin
CREATE TABLE tenants_properties (
    id          UUID PRIMARY KEY,
    tenant_id   UUID NOT NULL REFERENCES tenants(id)     ON DELETE CASCADE,
    property_id UUID NOT NULL REFERENCES properties(id)  ON DELETE CASCADE
);
CREATE INDEX idx_tenants_properties_tenant_id    ON tenants_properties(tenant_id);
CREATE INDEX idx_tenants_properties_property_id  ON tenants_properties(property_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tenants_properties;
-- +goose StatementEnd
