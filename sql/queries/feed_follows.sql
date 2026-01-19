-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4
    )
    RETURNING *
)
SELECT inserted_feed_follow.*, feeds.name as feed_name, users.name as user_name
FROM inserted_feed_follow 
INNER JOIN users ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name as feed_name, feeds.url, users.name as adding_user FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users ON feeds.user_id = users.id
WHERE $1 = feed_follows.user_id;

