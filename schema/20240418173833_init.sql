-- +goose Up
-- +goose StatementBegin
ALTER TABLE moderation ADD COLUMN content_hash varchar(64);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE moderation DROP COLUMN content_hash;
-- +goose StatementEnd
