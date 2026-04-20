-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS regulations (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    regulation_type VARCHAR(50) NOT NULL,
    issued_date DATE,
    status VARCHAR(50) NOT NULL DEFAULT 'Active'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS regulations;
-- +goose StatementEnd
