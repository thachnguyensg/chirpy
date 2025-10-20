-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (gen_random_uuid(), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $1, $2)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = $1;

-- name: DeleteChirp :one
DELETE FROM chirps WHERE id = $1
RETURNING *;
