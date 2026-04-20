-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS assessment_results (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES assessment_sessions(id) ON DELETE CASCADE,
    regulation_item_id UUID NOT NULL REFERENCES regulation_items(id) ON DELETE CASCADE,
    compliance_status VARCHAR(10) NOT NULL,
    evidence_link TEXT,
    remarks TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(session_id, regulation_item_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS assessment_results;
-- +goose StatementEnd
