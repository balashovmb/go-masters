-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
insert into users (name) values ('John'), ('Pit'), ('Bill');
insert into objects (name) values ('Iphone'), ('Macbook'), ('Ipad');

INSERT INTO reviews (user_id, object_id, text, rating) VALUES
(1, 1, 'Good phone', 2),
(1, 2, 'Bad pc', 1),
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
