-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

-- `пользователи`, `объекты`, `отзывы`
create table users (
    id serial primary key,
    name text 
);

create table objects (
    id serial primary key,
    name text
);

create table reviews (
    id serial primary key,
    user_id integer references users(id),
    object_id integer references objects(id),
    text text,
    rating integer
);

create INDEX idx_user_id ON reviews (user_id);
create INDEX idx_object_id ON reviews (object_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table users;
drop table objects;
drop table reviews;
-- +goose StatementEnd
