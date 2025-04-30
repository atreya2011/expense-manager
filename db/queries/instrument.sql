-- name: CreateInstrument :one
INSERT INTO
  instruments (id, name)
VALUES
  (lower(hex (randomblob (16))), ?) RETURNING *;

-- name: GetInstrument :one
SELECT
  *
FROM
  instruments
WHERE
  id = ?
LIMIT
  1;

-- name: ListInstruments :many
SELECT
  *
FROM
  instruments
ORDER BY
  name
LIMIT
  ?
OFFSET
  ?;

-- name: UpdateInstrument :one
UPDATE instruments
SET
  name = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE
  id = ? RETURNING *;

-- name: DeleteInstrument :exec
DELETE FROM instruments
WHERE
  id = ?;
