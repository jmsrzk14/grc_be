-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_categories DROP COLUMN tenant_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_categories ADD COLUMN tenant_id UUID;
CREATE INDEX idx_risk_categories_tenant_id ON risk_categories(tenant_id);
-- +goose StatementEnd
