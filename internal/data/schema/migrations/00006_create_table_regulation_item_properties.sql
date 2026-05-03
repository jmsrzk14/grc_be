-- +goose Up
-- +goose StatementBegin
CREATE TABLE regulation_item_properties (
    regulation_item_id UUID NOT NULL REFERENCES regulation_items(id) ON DELETE CASCADE,
    property_id        UUID NOT NULL REFERENCES properties(id)       ON DELETE CASCADE,
    PRIMARY KEY (regulation_item_id, property_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS regulation_item_properties;
-- +goose StatementEnd
