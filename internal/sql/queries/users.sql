-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE name = $1);

-- name: DeleteAll :exec
DELETE FROM users;

-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)

VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetUserFeeds :many
SELECT  users.name AS user_name, feeds.name AS feed_name, feeds.url AS feeds_url, feeds.id AS feed_id
FROM users
INNER JOIN feeds ON users.id = feeds.user_id;


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
