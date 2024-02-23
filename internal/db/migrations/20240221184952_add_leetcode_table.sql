-- +goose Up
-- +goose StatementBegin
CREATE TABLE leetcode
(
    user_id           INT REFERENCES users (id) UNIQUE,
    leetcode_username VARCHAR(50)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE leetcode;
-- +goose StatementEnd
