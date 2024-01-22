-- +goose Up
-- +goose StatementBegin
ALTER TABLE contributor
RENAME COLUMN username TO contributor
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
