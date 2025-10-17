-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, created_at, updated_at , expires_at, revoked_at)
VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $3, NULL)
RETURNING *;
