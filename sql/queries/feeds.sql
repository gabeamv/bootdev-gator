-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name, feeds.url, users.name as adding_user
FROM feeds INNER JOIN users ON feeds.user_id = users.id;

-- name: GetFeedFromUrl :one
SELECT * FROM feeds
WHERE $1 = feeds.url;