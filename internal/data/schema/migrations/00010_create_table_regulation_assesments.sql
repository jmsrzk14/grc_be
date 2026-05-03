-- +goose Up
-- +goose StatementBegin
CREATE TABLE regulation_assesments (
    id            UUID PRIMARY KEY,
    regulation_id UUID NOT NULL REFERENCES regulations(id)        ON DELETE CASCADE,
    session_id    UUID NOT NULL REFERENCES assessment_sessions(id) ON DELETE CASCADE,
    amount_pass   INT  DEFAULT 0,
    amount_fail   INT  DEFAULT 0,
    amount_na     INT  DEFAULT 0,
    UNIQUE (regulation_id, session_id)
);
CREATE INDEX idx_regulation_assesments_reg_id    ON regulation_assesments(regulation_id);
CREATE INDEX idx_regulation_assesments_ses_id    ON regulation_assesments(session_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS regulation_assesments;
-- +goose StatementEnd
