-- +goose Up
-- +goose StatementBegin
CREATE TABLE link_management (
    link uuid primary key,
    creator_id varchar not null,
    created_at timestamp not null
);
CREATE TABLE multiple (
    generated_link uuid references link_management(link) primary key ,
    creator_id varchar not null,
    racers varchar[] not null,
    track_name varchar not null default 'start_race',
    created_at timestamp not null,
    text_id uuid references text(id) not null
);
CREATE TABLE multiple_session(
    generated_link uuid references link_management(link),
    racer_id uuid references racer(id),
    duration int not null,
    wpm int not null,
    accuracy float not null,
    start_time timestamp not null ,
    winner varchar not null,
    place int not null,
    track_size int not null
);
CREATE TABLE multiple_history(
    generated_link uuid references link_management(link),
    racer_id uuid references racer(id),
    primary key(generated_link, racer_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE link_management;
DROP TABLE multiple_history;
DROP TABLE multiple_session;
DROP TABLE multiple;
-- +goose StatementEnd
