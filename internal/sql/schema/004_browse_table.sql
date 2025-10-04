-- +goose Up
CREATE TABLE browse (
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT,
    published_at TIMESTAMP,
    feed_id UUID NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS browse;
