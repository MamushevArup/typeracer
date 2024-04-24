-- +goose Up
-- +goose StatementBegin
create table avatar (
    id serial primary key,
    url varchar(512) not null
);
ALTER TABLE racer
    ADD COLUMN avatar_id INT,
    ADD CONSTRAINT fk_racer_avatar FOREIGN KEY (avatar_id) REFERENCES avatar(id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE racer DROP COLUMN avatar_id;
DROP TABLE avatar;
-- +goose StatementEnd
