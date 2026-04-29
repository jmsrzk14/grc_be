-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_categories (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    appetite TEXT,
    tolerance TEXT
);

CREATE TABLE risks (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    risk_title TEXT NOT NULL,
    risk_description TEXT,
    category_id UUID NOT NULL REFERENCES risk_categories(id),
    likelihood_inherent INTEGER DEFAULT 0,
    impact_inherent INTEGER DEFAULT 0,
    likelihood_residual INTEGER DEFAULT 0,
    impact_residual INTEGER DEFAULT 0,
    mitigation_plan TEXT,
    mitigation_status TEXT NOT NULL DEFAULT 'belum direncanakan'
);

CREATE INDEX idx_risk_categories_tenant_id ON risk_categories(tenant_id);
CREATE INDEX idx_risks_tenant_id ON risks(tenant_id);
CREATE INDEX idx_risks_category_id ON risks(category_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risks;
DROP TABLE risk_categories;
-- +goose StatementEnd
