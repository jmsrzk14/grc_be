-- +goose Up
-- +goose StatementBegin
CREATE TABLE risks (
    id                   UUID        PRIMARY KEY,
    risk_title           TEXT        NOT NULL,
    risk_description     TEXT,
    category_id          UUID        NOT NULL REFERENCES risk_categories(id) ON DELETE RESTRICT,
    likelihood_inherent  INTEGER     NOT NULL DEFAULT 0,
    impact_inherent      INTEGER     NOT NULL DEFAULT 0,
    likelihood_residual  INTEGER     NOT NULL DEFAULT 0,
    impact_residual      INTEGER     NOT NULL DEFAULT 0,
    mitigation_plan      TEXT,
    mitigation_status    TEXT        NOT NULL DEFAULT 'belum direncanakan'
);
CREATE INDEX idx_risks_category_id               ON risks(category_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risks;
-- +goose StatementEnd
