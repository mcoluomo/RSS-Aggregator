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
SELECT  users.name AS user_name, feeds.name AS feed_name, feeds.url AS feeds_url
FROM users
INNER JOIN feeds ON users.id = feeds.user_id;
