-- name: FollowFeed :exec
INSERT INTO feed_follow (id, user_id, feed_id)
VALUES ($1, $2, $3);