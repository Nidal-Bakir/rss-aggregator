-- name: FollowFeed :exec
INSERT INTO feed_follow (id, user_id, feed_id)
VALUES ($1, $2, $3);

-- name: GetFeedFollows :many
SELECT f.id,
    f.created_at,
    f.updated_at,
    f.name,
    f.url
FROM feed_follow AS ff
    JOIN feed As f ON f.id = ff.feed_id
WHERE ff.user_id = $1;