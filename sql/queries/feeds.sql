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

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $2, updated_at = $3
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;