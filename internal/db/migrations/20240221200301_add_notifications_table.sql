-- +goose Up
-- +goose StatementBegin
CREATE TABLE notifications
(
    user_id           INT REFERENCES users (id),
    notification_time TIME
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE notifications;
-- +goose StatementEnd
