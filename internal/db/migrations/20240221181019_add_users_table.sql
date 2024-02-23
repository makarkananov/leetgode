-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id            INT PRIMARY KEY,
    current_state VARCHAR(50)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
