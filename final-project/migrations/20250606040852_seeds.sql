-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
insert into users (name) values ('John'), ('Pit'), ('Bill');
insert into objects (name) values ('Iphone'), ('Macbook'), ('Ipad');
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
