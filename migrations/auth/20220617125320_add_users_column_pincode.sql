-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN pin_code varchar(6);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN pin_code;
-- +goose StatementEnd
