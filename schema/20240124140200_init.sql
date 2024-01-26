-- +goose Up
-- +goose StatementBegin
ALTER TABLE race_history
ADD COLUMN race_type int not null default 0
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
