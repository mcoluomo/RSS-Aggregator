
-- +goose Up
ALTER TABLE browse RENAME TO posts;

-- +goose Down
ALTER TABLE posts RENAME TO browse;
