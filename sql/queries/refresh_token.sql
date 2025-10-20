-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, created_at, updated_at , expires_at, revoked_at)
VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $3, NULL)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT u.*
FROM users u
JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE rt.token = $1 AND rt.revoked_at IS NULL AND rt.expires_at > CURRENT_TIMESTAMP;
