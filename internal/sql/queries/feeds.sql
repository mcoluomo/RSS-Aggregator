-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedId :one
SELECT feeds.id FROM feeds WHERE url = $1;
