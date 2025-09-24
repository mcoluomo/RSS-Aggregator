-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedId :one
SELECT feeds.id FROM feeds WHERE url = $1;
