-- +goose Up
-- +goose StatementBegin
ALTER TABLE contributor
ADD COLUMN username varchar(128);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
