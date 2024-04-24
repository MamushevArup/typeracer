-- +goose Up
-- +goose StatementBegin
create table avatar (
    id serial primary key,
    url varchar(512) not null
);
ALTER TABLE racer
    ADD COLUMN avatar_id INT,
    ADD CONSTRAINT fk_racer_avatar FOREIGN KEY (avatar_id) REFERENCES avatar(id);
ALTER TABLE racer
    DROP column avatar;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE racer DROP CONSTRAINT fk_racer_avatar;
ALTER TABLE racer DROP COLUMN avatar_id;
ALTER TABLE ADD COLUMN avatar varchar(512) not null;
DROP TABLE avatar;
-- +goose StatementEnd
