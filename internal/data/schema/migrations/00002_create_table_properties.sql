-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS properties (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS properties;
-- +goose StatementEnd
