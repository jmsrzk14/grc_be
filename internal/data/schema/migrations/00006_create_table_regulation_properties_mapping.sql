-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS regulation_properties_mapping (
    id UUID PRIMARY KEY,
    regulation_id UUID NOT NULL REFERENCES regulations(id) ON DELETE CASCADE,
    property_id UUID NOT NULL REFERENCES properties(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS regulation_properties_mapping;
-- +goose StatementEnd
