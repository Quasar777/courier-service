-- +goose Up
-- +goose StatementBegin
ALTER TABLE couriers
ALTER COLUMN status SET DEFAULT 'available';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE couriers
ALTER COLUMN status DROP DEFAULT;
-- +goose StatementEnd
