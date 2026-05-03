-- +goose Up
-- +goose StatementBegin
CREATE TABLE regulation_items (
    id               UUID PRIMARY KEY,
    regulation_id    UUID         NOT NULL REFERENCES regulations(id) ON DELETE CASCADE,
    reference_number VARCHAR(100) NOT NULL,
    content          TEXT
);
CREATE INDEX idx_regulation_items_regulation_id  ON regulation_items(regulation_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS regulation_items;
-- +goose StatementEnd
