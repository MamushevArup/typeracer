-- +goose Up
-- +goose StatementBegin
ALTER TABLE link_management
ADD COLUMN is_expired boolean not null default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE link_management
DROP COLUMN is_expired;
-- +goose StatementEnd
