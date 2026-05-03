-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_categories (
    id    UUID PRIMARY KEY,
    title TEXT NOT NULL,
    UNIQUE (title)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risk_categories;
-- +goose StatementEnd
