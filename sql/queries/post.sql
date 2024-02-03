-- name: CreatePost :exec
INSERT INTO post (id, title, url, pub_date, description, feed_id)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetPostsForFollowedFeed :many
SELECT post.*
FROM post
    JOIN feed_follow ON post.feed_id = feed_follow.feed_id
WHERE feed_follow.user_id = $1
ORDER BY pub_date DESC OFFSET $2
LIMIT $3;