-- name: CreateFeed :one
INSERT INTO feeds (
	id, 
	created_at,
	updated_at,
	name,
	url,
	user_id
)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)
RETURNING *;

-- name: GetAllFeeds :many
SELECT f.name, f.url, u.name AS user_name
FROM feeds f
JOIN users u
ON f.user_id = u.id;

-- name: GetFeedIdByUrl :one
SELECT id FROM feeds
WHERE url = $1;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $2, last_fetched_at = $2
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at
NULLS FIRST
LIMIT 1;
