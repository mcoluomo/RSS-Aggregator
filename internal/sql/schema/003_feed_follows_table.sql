-- +goose Up
CREATE TABLE feed_follows (
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    user_id UUID NOT NULL,
    feed_id UUID NOT NULL,
    PRIMARY KEY (user_id, feed_id),

        FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,

        FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE

);

-- +goose Down
DROP TABLE IF EXISTS feed_follows;
