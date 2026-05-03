-- +goose Up
-- +goose StatementBegin
CREATE TABLE assessment_results (
    id                 UUID        PRIMARY KEY,
    session_id         UUID        NOT NULL REFERENCES assessment_sessions(id) ON DELETE CASCADE,
    regulation_item_id UUID        NOT NULL REFERENCES regulation_items(id)    ON DELETE CASCADE,
    compliance_status  VARCHAR(10) NOT NULL,   -- YES | NO | N/A
    evidence_link      TEXT,
    remarks            TEXT,
    updated_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (session_id, regulation_item_id)
);
CREATE INDEX idx_assessment_results_session_id   ON assessment_results(session_id);
CREATE INDEX idx_assessment_results_item_id      ON assessment_results(regulation_item_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS assessment_results;
-- +goose StatementEnd
