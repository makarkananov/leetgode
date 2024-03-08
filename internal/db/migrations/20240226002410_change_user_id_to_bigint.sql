-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ALTER COLUMN id SET DATA TYPE BIGINT;

ALTER TABLE leetcode
    ALTER COLUMN user_id SET DATA TYPE BIGINT;

ALTER TABLE notifications
    ALTER COLUMN user_id SET DATA TYPE BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Если необходимо откатить изменения, можно восстановить предыдущий тип данных (INT):
ALTER TABLE users
    ALTER COLUMN id SET DATA TYPE INT;

ALTER TABLE leetcode
    ALTER COLUMN user_id SET DATA TYPE INT;

ALTER TABLE notifications
    ALTER COLUMN user_id SET DATA TYPE INT;
-- +goose StatementEnd
