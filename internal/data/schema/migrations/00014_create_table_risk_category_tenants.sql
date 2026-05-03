-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_category_tenants (
    id               UUID PRIMARY KEY,
    risk_id          UUID REFERENCES risks(id)           ON DELETE CASCADE,
    risk_category_id UUID NOT NULL REFERENCES risk_categories(id) ON DELETE CASCADE,
    tenant_id        UUID NOT NULL REFERENCES tenants(id)         ON DELETE CASCADE,
    appetite         TEXT,
    tolerance        TEXT,
    UNIQUE (risk_category_id, tenant_id)
);
CREATE INDEX idx_risk_category_tenants_cat_id    ON risk_category_tenants(risk_category_id);
CREATE INDEX idx_risk_category_tenants_tenant_id ON risk_category_tenants(tenant_id);
CREATE INDEX idx_risk_category_tenants_risk_id   ON risk_category_tenants(risk_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risk_category_tenants;
-- +goose StatementEnd
