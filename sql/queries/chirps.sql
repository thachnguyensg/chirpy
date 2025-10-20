-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (gen_random_uuid(), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $1, $2)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps;

-- name: GetChirpsByUserID :many
SELECT * FROM chirps WHERE user_id = $1;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = $1;

-- name: DeleteChirp :one
DELETE FROM chirps WHERE id = $1
RETURNING *;
