-- +goose Up
-- +goose StatementBegin
ALTER TABLE contributor DROP COLUMN username;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE contributor ADD COLUMN username VARCHAR(255) NOT NULL;
-- +goose StatementEnd
