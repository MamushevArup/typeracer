-- +goose Up
-- +goose StatementBegin
ALTER TABLE racer ADD COLUMN total_speed INT DEFAULT 0;
ALTER TABLE text ADD COLUMN total_speed INT DEFAULT 0;
ALTER TABLE multiple_session ALTER COLUMN accuracy TYPE int;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE racer DROP COLUMN total_speed;
ALTER TABLE multiple_session ALTER COLUMN accuracy TYPE float;
ALTER TABLE text DROP COLUMN total_speed;
-- +goose StatementEnd
