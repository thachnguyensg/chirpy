-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $1, $2)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUserAuth :one
UPDATE users
SET hashed_password = $2, email = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
