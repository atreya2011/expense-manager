-- name: CreateUser :one
INSERT INTO users (
  id, name, email
) VALUES (
  'usr_' || lower(hex(randomblob(16))), ?, ?
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: ListUsers :many
/* 
ListUsers returns a paginated list of users.
*/
SELECT * FROM users
ORDER BY name
LIMIT COALESCE(sqlc.narg('limit'), 50)
OFFSET COALESCE(sqlc.narg('offset'), 0);

-- name: UpdateUser :one
UPDATE users
SET name = ?, email = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
