-- name: CreatePost :one
INSERT INTO posts(
	id, 
	created_at, 
	updated_at,
	published_at,
	title,
	url,
	description,
	feed_id)
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsForUser :many
SELECT p.* FROM posts p
INNER JOIN feeds f 
ON p.feed_id = f.id
INNER JOIN users u
ON u.id = f.user_id
WHERE f.user_id = $1
ORDER BY published_at
LIMIT $2;
