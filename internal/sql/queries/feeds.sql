-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (user_id, feed_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4)
    RETURNING *
)
SELECT inserted_feed_follow.*,
    users.name AS user_name,
    feeds.name AS feed_name
FROM inserted_feed_follow
JOIN users ON inserted_feed_follow.user_id = users.id
JOIN feeds ON inserted_feed_follow.feed_id = users.id;


-- name: GetFeedId :one
SELECT feeds.id FROM feeds WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*,

    users.name AS user_name,
    feeds.name AS feed_name
    FROM feed_follows
JOIN users ON feed_follows.user_id = users.id
JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE users.id = $1;
