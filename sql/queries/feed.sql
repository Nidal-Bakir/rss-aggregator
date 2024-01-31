-- name: CreateFeed :one
INSERT INTO feed (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT id,
    created_at,
    updated_at,
    name,
    url
FROM feed;

-- name: GetFeedsOrderedByLastSync :many
SELECT id,
    url
FROM feed
ORDER BY last_sync_at ASC OFFSET $1
LIMIT $2;

-- name: MakeFeedAsSynced :exec
UPDATE feed
set last_sync_at = CURRENT_TIMESTAMP AT TIME ZONE 'UTC'
WHERE id = $1;